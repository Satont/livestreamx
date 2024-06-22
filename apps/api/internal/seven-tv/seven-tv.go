package seven_tv

import (
	"context"

	"github.com/satont/stream/apps/api/internal/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Config config.Config
}

func New(opts Opts) *SevenTV {
	s := &SevenTV{
		config: opts.Config,
		Emotes: make(map[string]Emote),
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go s.openWebSocket()
				go s.fetchEmotes()

				return nil
			},
		},
	)

	return s
}

type SevenTV struct {
	config config.Config

	Emotes map[string]Emote
}

type Emote struct {
	ID     string
	Name   string
	URL    string
	Width  int
	Height int
}
