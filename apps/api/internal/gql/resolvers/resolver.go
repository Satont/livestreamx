package resolvers

import (
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/satont/stream/apps/api/internal/config"
	"github.com/satont/stream/apps/api/internal/gql/mappers"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	mtx_api "github.com/satont/stream/apps/api/internal/mtx-api"
	chat_message "github.com/satont/stream/apps/api/internal/repositories/chat-message"
	message_reaction "github.com/satont/stream/apps/api/internal/repositories/message-reaction"
	"github.com/satont/stream/apps/api/internal/repositories/role"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	seven_tv "github.com/satont/stream/apps/api/internal/seven-tv"
	subscriptions_router "github.com/satont/stream/apps/api/internal/subscriptions-router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	chatMessageRepo     chat_message.Repository
	userRepo            user.Repository
	messageReactionRepo message_reaction.Repository
	rolesRepo           role.Repository
	// userFilesRepo            user_file.Repository

	sessionStorage     *session_storage.SessionStorage
	mapper             *mappers.Mapper
	s3                 *minio.Client
	config             config.Config
	sevenTv            *seven_tv.SevenTV
	mtxApi             *mtx_api.MtxApi
	subscriptionRouter subscriptions_router.Router
	redis              *redis.Client
	logger             *zap.Logger
}

type Opts struct {
	fx.In

	ChatMessageRepo     chat_message.Repository
	UserRepo            user.Repository
	MessageReactionRepo message_reaction.Repository
	RolesRepo           role.Repository
	// UserFilesRepo            user_file.Repository

	SessionStorage *session_storage.SessionStorage
	Converter      *mappers.Mapper
	// S3             *minio.Client
	Config             config.Config
	SevenTv            *seven_tv.SevenTV
	MtxApi             *mtx_api.MtxApi
	SubscriptionRouter subscriptions_router.Router
	Redis              *redis.Client
	Logger             *zap.Logger
}

func New(opts Opts) *Resolver {
	return &Resolver{
		chatMessageRepo:     opts.ChatMessageRepo,
		userRepo:            opts.UserRepo,
		messageReactionRepo: opts.MessageReactionRepo,
		rolesRepo:           opts.RolesRepo,
		subscriptionRouter:  opts.SubscriptionRouter,

		sessionStorage: opts.SessionStorage,
		mapper:         opts.Converter,
		config:         opts.Config,
		mtxApi:         opts.MtxApi,
		// userFilesRepo:            opts.UserFilesRepo,
		sevenTv: opts.SevenTv,
		redis:   opts.Redis,
		logger:  opts.Logger,
	}
}
