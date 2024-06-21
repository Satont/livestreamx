package converters

import (
	"github.com/satont/stream/apps/api/internal/repositories/user"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	UserRepo user.Repository
}

func New(opts Opts) *Converters {
	return &Converters{
		userRepo: opts.UserRepo,
	}
}

type Converters struct {
	userRepo user.Repository
}
