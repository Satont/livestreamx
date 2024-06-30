package httpserver

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/satont/stream/apps/api/internal/config"
	"github.com/satont/stream/apps/api/internal/gql"
	data_loader "github.com/satont/stream/apps/api/internal/gql/data-loader"
	"github.com/satont/stream/apps/api/internal/gql/mappers"
	"github.com/satont/stream/apps/api/internal/httpserver/middlewares"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC           fx.Lifecycle
	Gql          *gql.Gql
	Config       config.Config
	SessionStore *session_storage.SessionStorage
	Middlewares  *middlewares.Middlewares
	Mapper       *mappers.Mapper
	UserRepo     user.Repository
}

func New(opts Opts) *Server {
	gin.SetMode(gin.ReleaseMode)
	s := &Server{
		Engine: gin.New(),
	}

	s.Use(
		cors.New(
			cors.Config{
				AllowAllOrigins:  true,
				AllowMethods:     []string{"*"},
				AllowHeaders:     []string{"*"},
				ExposeHeaders:    []string{"*"},
				AllowCredentials: true,
			},
		),
		gin.Recovery(),
		sessions.Sessions("stream_session", opts.SessionStore),
		func(c *gin.Context) {
			c.Request = c.Request.WithContext(
				context.WithValue(c.Request.Context(), sessions.DefaultKey, sessions.Default(c)),
			)
			c.Next()
		},
		opts.Middlewares.AttachUserToContext,
	)

	s.GET(
		"/", func(c *gin.Context) {
			playground.Handler("GraphQL", "/api/query").ServeHTTP(c.Writer, c.Request)
		},
	)

	s.Any(
		"/query",
		func(c *gin.Context) {
			loader := data_loader.New(
				data_loader.Opts{
					UserRepo: opts.UserRepo,
					Mapper:   opts.Mapper,
				},
			)

			c.Request = c.Request.WithContext(
				context.WithValue(c.Request.Context(), data_loader.LoadersKey, loader),
			)
		},
		func(c *gin.Context) {
			opts.Gql.ServeHTTP(c.Writer, c.Request)
		},
	)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					s.StartServer(opts.Config.ApiPort)
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.StopServer()
				return nil
			},
		},
	)

	return s
}

type Server struct {
	*gin.Engine
}

func (c *Server) StartServer(port int) {
	c.Run(fmt.Sprintf(":%d", port))
}

func (c *Server) StopServer() {

}
