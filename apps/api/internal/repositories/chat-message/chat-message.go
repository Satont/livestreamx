package chat_message

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateChatMessageOpts) (*Message, error)
	FindLatest(ctx context.Context, opts FindLatestOpts) ([]Message, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Message, error)
}

type CreateChatMessageOpts struct {
	ChannelID uuid.UUID
	SenderID  uuid.UUID
	Text      string
	ReplyTo   *uuid.UUID
}

type FindLatestOpts struct {
	ChannelID uuid.UUID
	Limit     int
}
