package seven_tv

import (
	"github.com/satont/stream/apps/api/internal/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config config.Config
}

func New(opts Opts) *SevenTV {
	NewWs(opts.Config.SevenTVEmoteSetID)

	return &SevenTV{}
}

type SevenTV struct {
	emotes map[string]Emote
}

type Emote struct {
	ID   string
	Name string
	URL  string
}
