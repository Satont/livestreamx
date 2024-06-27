package mappers

import (
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

func (c *Mapper) DbProviderToGql(p user.UserConnectionProvider) gqlmodel.AuthedUserProviderType {
	switch p {
	case user.UserConnectionProviderTwitch:
		return gqlmodel.AuthedUserProviderTypeTwitch
	case user.UserConnectionProviderGithub:
		return gqlmodel.AuthedUserProviderTypeGithub
	default:
		return gqlmodel.AuthedUserProviderTypeTwitch
	}
}
