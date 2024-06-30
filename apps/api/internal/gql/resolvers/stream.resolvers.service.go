package resolvers

import (
	"fmt"

	"github.com/google/uuid"
)

func (r *Resolver) computeStreamThumbnailUrl(channelID uuid.UUID) string {
	return fmt.Sprintf("%s/%s", r.config.ThumbnailsURI, channelID.String())
}
