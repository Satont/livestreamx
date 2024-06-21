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
	"github.com/satont/stream/apps/api/internal/gql/converters"
	"github.com/satont/stream/apps/api/internal/gql/directives"
	"github.com/satont/stream/apps/api/internal/gql/resolvers"
	"github.com/satont/stream/apps/api/internal/httpserver"
	"github.com/satont/stream/apps/api/internal/httpserver/middlewares"
	"github.com/satont/stream/apps/api/internal/httpserver/routes/auth"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	chat_message "github.com/satont/stream/apps/api/internal/repositories/chat-message"
	chat_messages_with_user "github.com/satont/stream/apps/api/internal/repositories/chat-messages-with-user"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	seven_tv "github.com/satont/stream/apps/api/internal/seven-tv"
	"go.uber.org/fx"
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
	),

	fx.Provide(
		fx.Annotate(chat_message.NewPgx, fx.As(new(chat_message.Repository))),
		fx.Annotate(user.NewPgx, fx.As(new(user.Repository))),
		fx.Annotate(chat_messages_with_user.NewPgx, fx.As(new(chat_messages_with_user.Repository))),
	),

	fx.Provide(
		converters.New,
		session_storage.New,
		resolvers.New,
		middlewares.New,
		directives.New,
		gql.New,
		httpserver.New,
	),
	fx.Invoke(
		func(p *pgxpool.Pool, c config.Config) error {
			go seven_tv.NewWs(c.SevenTVEmoteSetID)
			if err := goose.SetDialect("pgx"); err != nil {
				return err
			}
			db, err := sql.Open("pgx", c.PostgresURL)
			if err != nil {
				return err
			}

			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			return goose.Up(db, filepath.Join(wd, "migrations"))
		},
		auth.New,
	),
)
