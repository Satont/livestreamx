package main

import (
	"github.com/satont/stream/apps/api/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.App,
	).Run()
}
