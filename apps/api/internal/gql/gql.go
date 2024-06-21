package gql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/satont/stream/apps/api/internal/gql/directives"
	"github.com/satont/stream/apps/api/internal/gql/graph"
	"github.com/satont/stream/apps/api/internal/gql/resolvers"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	"go.uber.org/fx"
)

type Gql struct {
	*handler.Server
}

type Opts struct {
	fx.In

	Resolver       *resolvers.Resolver
	SessionStorage *session_storage.SessionStorage
	Directives     *directives.Directives
}

func New(opts Opts) *Gql {
	graphConfig := graph.Config{
		Resolvers: opts.Resolver,
	}
	graphConfig.Directives.IsAuthenticated = opts.Directives.IsAuthenticated
	graphConfig.Directives.NotBanned = opts.Directives.NotBanned
	graphConfig.Directives.HasFeature = opts.Directives.HasFeature

	schema := graph.NewExecutableSchema(graphConfig)

	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(
		transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second,
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	)

	srv.Use(extension.Introspection{})

	return &Gql{srv}
}
