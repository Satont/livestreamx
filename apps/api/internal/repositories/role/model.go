package role

import (
	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID
	ChannelID uuid.UUID
	Name      string
	ImageUrl  *string
	Features  []string
}
