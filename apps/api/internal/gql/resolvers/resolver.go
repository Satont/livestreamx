package resolvers

import (
	"github.com/satont/stream/apps/api/internal/gql/converters"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	chat_message "github.com/satont/stream/apps/api/internal/repositories/chat-message"
	chat_messages_with_user "github.com/satont/stream/apps/api/internal/repositories/chat-messages-with-user"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	"go.uber.org/fx"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	chatMessageRepo          chat_message.Repository
	userRepo                 user.Repository
	chatMessagesWithUserRepo chat_messages_with_user.Repository
	sessionStorage           *session_storage.SessionStorage
	converter                *converters.Converters

	chatListenersChannels map[string]chan *gqlmodel.ChatMessage
}

type Opts struct {
	fx.In

	ChatMessageRepo          chat_message.Repository
	UserRepo                 user.Repository
	SessionStorage           *session_storage.SessionStorage
	ChatMessagesWithUserRepo chat_messages_with_user.Repository
	Converter                *converters.Converters
}

func New(opts Opts) *Resolver {
	return &Resolver{
		chatMessageRepo:          opts.ChatMessageRepo,
		userRepo:                 opts.UserRepo,
		sessionStorage:           opts.SessionStorage,
		chatMessagesWithUserRepo: opts.ChatMessagesWithUserRepo,
		chatListenersChannels:    make(map[string]chan *gqlmodel.ChatMessage),
		converter:                opts.Converter,
	}
}
