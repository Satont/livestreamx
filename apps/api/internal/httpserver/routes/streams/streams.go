package streams

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/stream/apps/api/internal/config"
	"github.com/satont/stream/apps/api/internal/httpserver"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	HttpServer *httpserver.Server
	Config     config.Config
	UserRepo   user.Repository
	Redis      *redis.Client
}

func New(opts Opts) (*Streams, error) {
	s := &Streams{
		config:   opts.Config,
		userRepo: opts.UserRepo,
		redis:    opts.Redis,
	}

	group := opts.HttpServer.Group("/streams")
	group.POST("auth", s.authHandler)

	group.GET("/*regex", s.reverseProxy("http://localhost:8888"))

	return s, nil
}

type Streams struct {
	config   config.Config
	userRepo user.Repository
	redis    *redis.Client
}
