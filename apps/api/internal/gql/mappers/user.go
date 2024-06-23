package mappers

import (
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

func (c *Converters) DbUserToGql(u user.User) gqlmodel.User {
	return gqlmodel.User{
		ID:          u.ID.String(),
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Color:       u.Color,
		Roles:       nil,
		IsBanned:    u.Banned,
		CreatedAt:   u.CreatedAt,
		AvatarURL:   u.AvatarUrl,
	}
}
