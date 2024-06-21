package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/satont/stream/apps/api/internal/httpserver/middlewares"
)

func (c *Directives) NotBanned(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (res interface{}, err error) {
	user := middlewares.GetUserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unauthenticated")
	}

	if user.Banned {
		return nil, fmt.Errorf("user is banned")
	}

	return next(ctx)
}
