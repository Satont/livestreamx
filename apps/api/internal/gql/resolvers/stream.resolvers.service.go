package resolvers

import (
	"fmt"

	"github.com/google/uuid"
)

func (r *Resolver) computeStreamThumbnailUrl(channelID uuid.UUID) string {
	return fmt.Sprintf("%s/%s", r.config.ThumbnailsURI, channelID.String())
}

func (r *Resolver) buildStreamViewersRedisKey(channelID uuid.UUID) string {
	return fmt.Sprintf("stream:%s:viewers", channelID.String())
}

func (r *Resolver) buildStreamChattersRedisKey(channelID uuid.UUID) string {
	return fmt.Sprintf("streams:chatters:%s", channelID.String())
}
