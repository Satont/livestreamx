package chat_messages_with_user

import (
	chat_message "github.com/satont/stream/apps/api/internal/repositories/chat-message"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

type MessageWithUser struct {
	Message chat_message.Message
	User    user.User
}
