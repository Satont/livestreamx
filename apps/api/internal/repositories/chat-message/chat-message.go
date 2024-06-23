package chat_message

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateChatMessageOpts) (*Message, error)
	FindLatest(ctx context.Context, opts FindManyOpts) ([]Message, error)
}

type CreateChatMessageOpts struct {
	SenderID uuid.UUID
	Text     string
	ReplyTo  *uuid.UUID
}

type FindManyOpts struct {
	Limit int
}
