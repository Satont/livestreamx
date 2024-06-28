package app

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"github.com/satont/stream/apps/api/internal/config"
	"github.com/satont/stream/apps/api/internal/gql"
	"github.com/satont/stream/apps/api/internal/gql/directives"
	"github.com/satont/stream/apps/api/internal/gql/mappers"
	"github.com/satont/stream/apps/api/internal/gql/resolvers"
	"github.com/satont/stream/apps/api/internal/httpserver"
	"github.com/satont/stream/apps/api/internal/httpserver/middlewares"
	"github.com/satont/stream/apps/api/internal/httpserver/routes/auth"
	"github.com/satont/stream/apps/api/internal/httpserver/routes/streams"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	mtx_api "github.com/satont/stream/apps/api/internal/mtx-api"
	chat_message "github.com/satont/stream/apps/api/internal/repositories/chat-message"
	chat_messages_with_user "github.com/satont/stream/apps/api/internal/repositories/chat-messages-with-user"
	message_reaction "github.com/satont/stream/apps/api/internal/repositories/message-reaction"
	"github.com/satont/stream/apps/api/internal/repositories/role"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	seven_tv "github.com/satont/stream/apps/api/internal/seven-tv"
	subscriptions_router "github.com/satont/stream/apps/api/internal/subscriptions-router"
	"go.uber.org/fx"

	"github.com/nats-io/nats.go"
)

var App = fx.Options(
	fx.Provide(
		config.New,
		func(c config.Config) (*redis.Client, error) {
			opts, err := redis.ParseURL(c.RedisURL)
			if err != nil {
				return nil, err
			}

			return redis.NewClient(opts), nil
		},
		func(c config.Config) (*pgxpool.Pool, error) {
			return pgxpool.New(context.Background(), c.PostgresURL)
		},
		func(c config.Config) (*nats.Conn, error) {
			nc, err := nats.Connect(c.NatsURL)
			if err != nil {
				return nil, err
			}

			return nc, nil
		},
	),

	fx.Provide(
		fx.Annotate(chat_message.NewPgx, fx.As(new(chat_message.Repository))),
		fx.Annotate(user.NewPgx, fx.As(new(user.Repository))),
		fx.Annotate(chat_messages_with_user.NewPgx, fx.As(new(chat_messages_with_user.Repository))),
		fx.Annotate(message_reaction.NewPgx, fx.As(new(message_reaction.Repository))),
		fx.Annotate(role.NewPgx, fx.As(new(role.Repository))),
	),

	fx.Provide(
		// s3.New,
		fx.Annotate(
			subscriptions_router.NewNatsSubscription,
			fx.As(new(subscriptions_router.Router)),
		),
		mtx_api.New,
		seven_tv.New,
		mappers.New,
		session_storage.New,
		resolvers.New,
		middlewares.New,
		directives.New,
		gql.New,
		httpserver.New,
	),
	fx.Invoke(
		func(p *pgxpool.Pool, c config.Config) error {
			if err := goose.SetDialect("pgx"); err != nil {
				return err
			}
			db, err := sql.Open("pgx", c.PostgresURL)
			if err != nil {
				return err
			}
			defer db.Close()

			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			return goose.Up(db, filepath.Join(wd, "migrations"))
		},
		auth.New,
		streams.New,
	),
)
