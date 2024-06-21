package resolvers

import (
	"github.com/minio/minio-go/v7"
	"github.com/satont/stream/apps/api/internal/config"
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
	// userFilesRepo            user_file.Repository

	sessionStorage *session_storage.SessionStorage
	converter      *converters.Converters
	s3             *minio.Client
	config         config.Config

	chatListenersChannels map[string]chan *gqlmodel.ChatMessage
	streamViewers         int
	streamChatters        map[string]gqlmodel.Chatter
}

type Opts struct {
	fx.In

	ChatMessageRepo          chat_message.Repository
	UserRepo                 user.Repository
	ChatMessagesWithUserRepo chat_messages_with_user.Repository
	// UserFilesRepo            user_file.Repository

	SessionStorage *session_storage.SessionStorage
	Converter      *converters.Converters
	// S3             *minio.Client
	Config config.Config
}

func New(opts Opts) *Resolver {
	return &Resolver{
		chatMessageRepo:          opts.ChatMessageRepo,
		userRepo:                 opts.UserRepo,
		chatMessagesWithUserRepo: opts.ChatMessagesWithUserRepo,
		sessionStorage:           opts.SessionStorage,
		converter:                opts.Converter,
		config:                   config.Config{},
		// userFilesRepo:            opts.UserFilesRepo,
		chatListenersChannels: make(map[string]chan *gqlmodel.ChatMessage),
		streamViewers:         0,
		streamChatters:        make(map[string]gqlmodel.Chatter),
	}
}
