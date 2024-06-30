package chat_message

import (
	"time"

	"github.com/google/uuid"
	message_reaction "github.com/satont/stream/apps/api/internal/repositories/message-reaction"
)

type Message struct {
	ID        uuid.UUID
	ChannelID uuid.UUID
	SenderID  uuid.UUID
	Text      string
	CreatedAt time.Time
	ReplyTo   *uuid.UUID

	Reactions []message_reaction.MessageReaction
}
