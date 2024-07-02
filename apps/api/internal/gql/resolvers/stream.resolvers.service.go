package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Resolver) resetStreamsViewersAndChatters(ctx context.Context) error {
	streamsKey := "streams:viewers:*"
	chattersKey := "streams:chatters:*"

	for _, key := range []string{streamsKey, chattersKey} {
		var cursor uint64
		var keys []string
		var err error
		keys, cursor, err = r.redis.Scan(ctx, cursor, key, 0).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		for _, key := range keys {
			if err := r.redis.Del(ctx, key).Err(); err != nil {
				return fmt.Errorf("failed to delete key: %w", err)
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
