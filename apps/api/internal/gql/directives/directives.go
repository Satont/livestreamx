package directives

import (
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	SessionStorage *session_storage.SessionStorage
	UserRepo       user.Repository
}

func New(opts Opts) *Directives {
	return &Directives{
		sessionStorage: opts.SessionStorage,
		userRepo:       opts.UserRepo,
	}
}

type Directives struct {
	sessionStorage *session_storage.SessionStorage
	userRepo       user.Repository
}
