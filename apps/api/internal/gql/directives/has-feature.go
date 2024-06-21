package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
)

func (c *Directives) HasFeature(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
	features []gqlmodel.RoleFeature,
) (res interface{}, err error) {
	return next(ctx)
}
