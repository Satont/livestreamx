package streams

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/stream/apps/api/internal/config"
	"github.com/satont/stream/apps/api/internal/httpserver"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Opts struct {
	fx.In

	HttpServer *httpserver.Server
	Config     config.Config
	UserRepo   user.Repository
	Redis      *redis.Client
	Logger     *zap.Logger
}

func New(opts Opts) (*Streams, error) {
	s := &Streams{
		config:   opts.Config,
		userRepo: opts.UserRepo,
		redis:    opts.Redis,
		logger:   opts.Logger,
	}

	group := opts.HttpServer.Group("/streams")
	group.POST("auth", s.authHandler)
	group.GET("/thumbnails/*channelID", s.thumbnailsHandler)

	group.GET("/:channelID/index.m3u8", s.indexHandler)
	group.Any("/read/*regex", s.reverseProxy(opts.Config.MediaMtxAddr+":8888"))

	return s, nil
}

type Streams struct {
	config   config.Config
	userRepo user.Repository
	redis    *redis.Client
	logger   *zap.Logger
}
