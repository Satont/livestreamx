package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Resolver) resetStreamsViewersAndChatters(ctx context.Context) error {
	streamsKey := "streams:viewers:*"
	chattersKey := "streams:chatters:*"

	for _, keyPattern := range []string{streamsKey, chattersKey} {
		iter := r.redis.Scan(ctx, 0, keyPattern, 0).Iterator()

		for iter.Next(ctx) {
			err := r.redis.Del(ctx, iter.Val()).Err()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Resolver) computeStreamThumbnailUrl(channelID uuid.UUID) string {
	return fmt.Sprintf("%s/%s", r.config.ThumbnailsURI, channelID.String())
}

func (r *Resolver) buildStreamViewersRedisKey(channelID uuid.UUID) string {
	return fmt.Sprintf("streams:viewers:%s", channelID.String())
}

func (r *Resolver) buildStreamChattersRedisKey(channelID uuid.UUID) string {
	return fmt.Sprintf("streams:chatters:%s", channelID.String())
}
