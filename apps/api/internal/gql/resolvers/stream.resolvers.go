package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	data_loader "github.com/satont/stream/apps/api/internal/gql/data-loader"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/gql/graph"
	"github.com/satont/stream/apps/api/internal/httpserver/middlewares"
	userrepo "github.com/satont/stream/apps/api/internal/repositories/user"
	"golang.org/x/sync/errgroup"
)

// Streams is the resolver for the streams field.
func (r *queryResolver) Streams(ctx context.Context) ([]gqlmodel.Stream, error) {
	paths, err := r.mtxApi.GetPaths(ctx)
	if err != nil {
		return nil, err
	}

	var streamsMu sync.Mutex
	var streams []gqlmodel.Stream
	var errwg errgroup.Group

	for _, path := range paths {
		path := path
		errwg.Go(
			func() error {
				parsedStreamKey, err := uuid.Parse(path.Name)
				if err != nil {
					return err
				}

				dbChannel, err := r.userRepo.FindByStreamKey(ctx, parsedStreamKey)
				if err != nil {
					return err
				}

				streamsMu.Lock()
				defer streamsMu.Unlock()

				streams = append(
					streams,
					gqlmodel.Stream{
						Chatters:     []gqlmodel.Chatter{},
						StartedAt:    path.ReadyTime,
						ChannelID:    dbChannel.ID,
						ThumbnailURL: r.Resolver.computeStreamThumbnailUrl(dbChannel.ID),
					},
				)

				return nil
			},
		)
	}

	if err := errwg.Wait(); err != nil {
		return nil, err
	}

	return streams, nil
}

// Viewers is the resolver for the viewers field.
func (r *streamResolver) Viewers(ctx context.Context, obj *gqlmodel.Stream) (int, error) {
	viewers, err := r.redis.Get(ctx, r.buildStreamViewersRedisKey(obj.ChannelID)).Int()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			r.logger.Sugar().Error(err)
		}
		return 0, nil
	}

	return viewers, nil
}

// Chatters is the resolver for the chatters field.
func (r *streamResolver) Chatters(ctx context.Context, obj *gqlmodel.Stream) (
	[]gqlmodel.Chatter,
	error,
) {
	var chatters []gqlmodel.Chatter
	iter := r.redis.Scan(
		ctx,
		0,
		fmt.Sprintf("%s:*", r.buildStreamChattersRedisKey(obj.ChannelID)),
		0,
	).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := r.redis.Get(ctx, key).Bytes()
		if err != nil {
			r.logger.Sugar().Error(err)
			return nil, err
		}

		var u userrepo.User
		if err := json.Unmarshal(val, &u); err != nil {
			r.logger.Sugar().Error(err)
			return nil, err
		}

		gqlUser := r.mapper.DbUserToBaseUserGql(u)

		chatters = append(
			chatters,
			gqlmodel.Chatter{
				User: &gqlUser,
			},
		)
	}
	if err := iter.Err(); err != nil {
		r.logger.Sugar().Error(err)
		return nil, err
	}

	return chatters, nil
}

// Channel is the resolver for the channel field.
func (r *streamResolver) Channel(ctx context.Context, obj *gqlmodel.Stream) (
	*gqlmodel.BaseUser,
	error,
) {
	return data_loader.GetBaseUserByID(ctx, obj.ChannelID)
}

// StreamInfo is the resolver for the streamInfo field.
func (r *subscriptionResolver) StreamInfo(
	ctx context.Context,
	channelID uuid.UUID,
) (<-chan *gqlmodel.Stream, error) {
	dbChannel, err := r.userRepo.FindByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	user := middlewares.GetUserFromContext(ctx)
	if user != nil {
		userBytes, err := json.Marshal(user)
		if err != nil {
			r.logger.Sugar().Error(err)
			return nil, err
		}
		userChatterKey := fmt.Sprintf("%s:%s", r.buildStreamChattersRedisKey(channelID), user.ID)
		if err := r.redis.Set(ctx, userChatterKey, userBytes, 48*time.Hour).Err(); err != nil {
			return nil, err
		}
	}

	channel := make(chan *gqlmodel.Stream, 1)

	if err := r.redis.Incr(ctx, r.buildStreamViewersRedisKey(channelID)).Err(); err != nil {
		return nil, err
	}

	if err := r.redis.Expire(
		ctx,
		r.buildStreamViewersRedisKey(channelID),
		48*time.Hour,
	).Err(); err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			close(channel)
			if user != nil {
				userChatterKey := fmt.Sprintf("%s:%s", r.buildStreamChattersRedisKey(channelID), user.ID)
				if err := r.redis.Del(context.Background(), userChatterKey).Err(); err != nil {
					r.logger.Sugar().Error(err)
				}
			}

			if err := r.redis.Decr(
				context.Background(),
				r.buildStreamViewersRedisKey(channelID),
			).Err(); err != nil {
				r.logger.Sugar().Error(err)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				mtxInfo, err := r.mtxApi.GetPathInfoByApyKey(ctx, dbChannel.StreamKey.String())
				if err != nil {
					r.logger.Sugar().Error(err)
					time.Sleep(1 * time.Second)
					continue
				}

				streamInfo := &gqlmodel.Stream{
					StartedAt:    mtxInfo.ReadyTime,
					ChannelID:    dbChannel.ID,
					ThumbnailURL: r.Resolver.computeStreamThumbnailUrl(dbChannel.ID),
				}

				channel <- streamInfo
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return channel, nil
}

// Stream returns graph.StreamResolver implementation.
func (r *Resolver) Stream() graph.StreamResolver { return &streamResolver{r} }

type streamResolver struct{ *Resolver }
