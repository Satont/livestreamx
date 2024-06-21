package chat_message

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	SenderID  uuid.UUID
	Text      string
	CreatedAt time.Time
}
