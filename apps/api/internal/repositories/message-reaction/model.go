package message_reaction

import (
	"time"

	"github.com/google/uuid"
)

type MessageReaction struct {
	ID        uuid.UUID
	MessageID uuid.UUID
	UserID    uuid.UUID
	Reaction  string
	CreatedAt time.Time
}
