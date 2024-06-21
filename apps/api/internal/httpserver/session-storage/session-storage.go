package session_storage

import (
	"context"
	"errors"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/satont/stream/apps/api/internal/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config config.Config
}

const USER_ID_KEY = "user_id"

func New(opts Opts) (*SessionStorage, error) {
	redisOpts, err := goredis.ParseURL(opts.Config.RedisURL)
	if err != nil {
		return nil, err
	}

	sessionAge := (31 * 24 * time.Hour).Seconds()
	sessionAgeSeconds := int(sessionAge)

	store, _ := redis.NewStore(
		10,
		"tcp",
		redisOpts.Addr,
		redisOpts.Password,
		[]byte(opts.Config.ApiSessionSecret),
	)
	store.Options(
		sessions.Options{
			MaxAge:   sessionAgeSeconds,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		},
	)

	return &SessionStorage{
		store,
	}, nil
}

type SessionStorage struct {
	sessions.Store
}

func (c *SessionStorage) GetSession(ctx context.Context) sessions.Session {
	return ctx.Value(sessions.DefaultKey).(sessions.Session)
}

func (c *SessionStorage) GetUserID(ctx context.Context) (string, error) {
	session := c.GetSession(ctx)
	userID := session.Get(USER_ID_KEY)
	if userID == nil {
		return "", errors.New("user_id not found in session")
	}
	casted, ok := userID.(string)
	if !ok {
		return "", errors.New("user_id not a string")
	}

	return casted, nil
}
