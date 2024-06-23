package message_reaction

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindOne(ctx context.Context, reactionID uuid.UUID) (
		*MessageReaction,
		error,
	)
	FindManyByMessageID(ctx context.Context, messageID uuid.UUID) ([]MessageReaction, error)
	Create(ctx context.Context, opts CreateMessageReactionOpts) (*MessageReaction, error)
	Remove(ctx context.Context, id uuid.UUID) error
}

type CreateMessageReactionOpts struct {
	MessageID uuid.UUID
	UserID    uuid.UUID
	Reaction  string
}
