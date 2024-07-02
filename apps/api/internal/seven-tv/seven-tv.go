package seven_tv

import (
	"context"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/satont/stream/apps/api/internal/config"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	subscriptions_router "github.com/satont/stream/apps/api/internal/subscriptions-router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Config              config.Config
	UserRepo            user.Repository
	Logger              *zap.Logger
	SubscriptionsRouter subscriptions_router.Router
}

func New(opts Opts) *SevenTV {
	s := &SevenTV{
		config:              opts.Config,
		Channels:            make([]ChannelCache, 0),
		userRepo:            opts.UserRepo,
		logger:              opts.Logger,
		subscriptionsRouter: opts.SubscriptionsRouter,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go s.openWebSocket()
				go s.init()

				return nil
			},
		},
	)

	return s
}

type SevenTV struct {
	config              config.Config
	userRepo            user.Repository
	wsConn              *websocket.Conn
	logger              *zap.Logger
	subscriptionsRouter subscriptions_router.Router

	// map of channels with emotes
	Channels []ChannelCache
}

type ChannelCache struct {
	EmoteSetID string
	ChannelID  uuid.UUID
	Emotes     map[string]Emote
}

type Emote struct {
	ID     string
	Name   string
	URL    string
	Width  int
	Height int
}
