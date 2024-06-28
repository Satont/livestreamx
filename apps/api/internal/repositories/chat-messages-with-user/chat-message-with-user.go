package chat_messages_with_user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindLatest(ctx context.Context, opts FindLatestOpts) ([]MessageWithUser, error)
	FindByID(ctx context.Context, id uuid.UUID) (*MessageWithUser, error)
}

type FindLatestOpts struct {
	Limit     int
	ChannelID uuid.UUID
}
