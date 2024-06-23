package resolvers

import (
	"github.com/satont/stream/apps/api/internal/gql/mappers"
	mtx_api "github.com/satont/stream/apps/api/internal/mtx-api"
	message_reaction "github.com/satont/stream/apps/api/internal/repositories/message-reaction"
	"go.uber.org/atomic"

	"github.com/minio/minio-go/v7"
	"github.com/satont/stream/apps/api/internal/config"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	chat_message "github.com/satont/stream/apps/api/internal/repositories/chat-message"
	chat_messages_with_user "github.com/satont/stream/apps/api/internal/repositories/chat-messages-with-user"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	seven_tv "github.com/satont/stream/apps/api/internal/seven-tv"
	"go.uber.org/fx"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	chatMessageRepo          chat_message.Repository
	userRepo                 user.Repository
	chatMessagesWithUserRepo chat_messages_with_user.Repository
	messageReactionRepo      message_reaction.Repository
	// userFilesRepo            user_file.Repository

	sessionStorage *session_storage.SessionStorage
	mapper         *mappers.Converters
	s3             *minio.Client
	config         config.Config
	sevenTv        *seven_tv.SevenTV
	mtxApi         *mtx_api.MtxApi

	chatListenersChannels     map[string]chan *gqlmodel.ChatMessage
	reactionListenersChannels map[string]chan gqlmodel.ChatMessageReaction
	streamViewers             *atomic.Int32
	streamChatters            map[string]gqlmodel.Chatter
}

type Opts struct {
	fx.In

	ChatMessageRepo          chat_message.Repository
	UserRepo                 user.Repository
	ChatMessagesWithUserRepo chat_messages_with_user.Repository
	MessageReactionRepo      message_reaction.Repository
	// UserFilesRepo            user_file.Repository

	SessionStorage *session_storage.SessionStorage
	Converter      *mappers.Converters
	// S3             *minio.Client
	Config  config.Config
	SevenTv *seven_tv.SevenTV
	MtxApi  *mtx_api.MtxApi
}

func New(opts Opts) *Resolver {
	return &Resolver{
		chatMessageRepo:          opts.ChatMessageRepo,
		userRepo:                 opts.UserRepo,
		chatMessagesWithUserRepo: opts.ChatMessagesWithUserRepo,
		messageReactionRepo:      opts.MessageReactionRepo,

		sessionStorage: opts.SessionStorage,
		mapper:         opts.Converter,
		config:         opts.Config,
		mtxApi:         opts.MtxApi,
		// userFilesRepo:            opts.UserFilesRepo,
		sevenTv:                   opts.SevenTv,
		chatListenersChannels:     make(map[string]chan *gqlmodel.ChatMessage),
		reactionListenersChannels: make(map[string]chan gqlmodel.ChatMessageReaction),
		streamViewers:             atomic.NewInt32(0),
		streamChatters:            make(map[string]gqlmodel.Chatter),
	}
}
