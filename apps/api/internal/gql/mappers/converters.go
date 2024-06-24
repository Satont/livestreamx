package mappers

import (
	"github.com/satont/stream/apps/api/internal/repositories/user"
	seven_tv "github.com/satont/stream/apps/api/internal/seven-tv"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	UserRepo user.Repository
	SevenTv  *seven_tv.SevenTV
}

func New(opts Opts) *Mapper {
	return &Mapper{
		userRepo: opts.UserRepo,
		sevenTv:  opts.SevenTv,
	}
}

type Mapper struct {
	userRepo user.Repository
	sevenTv  *seven_tv.SevenTV
}
