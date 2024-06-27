package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/satont/stream/apps/api/internal/config"
	"github.com/satont/stream/apps/api/internal/httpserver"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	HttpServer   *httpserver.Server
	Config       config.Config
	SessionStore *session_storage.SessionStorage
	UserRepo     user.Repository
}

func New(opts Opts) (*Auth, error) {
	a := &Auth{
		config:       opts.Config,
		sessionStore: opts.SessionStore,
		userRepo:     opts.UserRepo,
	}

	group := opts.HttpServer.Group("/auth")
	group.GET(
		"/", func(context *gin.Context) {
			s := opts.SessionStore.GetSession(context.Request.Context())

			context.JSON(200, gin.H{"user_id": s.Get("user_id")})
		},
	)

	group.GET("/twitch", a.TwitchGetLink)
	group.GET("/twitch/callback", a.TwitchCallback)
	group.GET("/github", a.GithubGetLink)
	group.GET("/github/callback", a.GithubCallback)

	return a, nil
}

type Auth struct {
	config       config.Config
	sessionStore *session_storage.SessionStorage
	userRepo     user.Repository
}
