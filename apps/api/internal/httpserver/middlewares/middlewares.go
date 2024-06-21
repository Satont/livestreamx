package middlewares

import (
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	UserRepo       user.Repository
	SessionStorage *session_storage.SessionStorage
}

func New(opts Opts) *Middlewares {
	return &Middlewares{
		userRepo:       opts.UserRepo,
		sessionStorage: opts.SessionStorage,
	}
}

type Middlewares struct {
	userRepo       user.Repository
	sessionStorage *session_storage.SessionStorage
}
