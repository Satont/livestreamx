package mappers

import (
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

func (c *Mapper) DbUserToBaseUserGql(u user.User) gqlmodel.BaseUser {
	return gqlmodel.BaseUser{
		ID:          u.ID,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Color:       u.Color,
		IsBanned:    u.Banned,
		CreatedAt:   u.CreatedAt,
		AvatarURL:   u.AvatarUrl,
		IsAdmin:     u.IsAdmin,
	}
}

func (c *Mapper) DbUserToChatUser(u user.User) gqlmodel.ChatUser {
	return gqlmodel.ChatUser{
		ID:          u.ID,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Color:       u.Color,
		IsBanned:    u.Banned,
		CreatedAt:   u.CreatedAt,
		AvatarURL:   u.AvatarUrl,
		IsAdmin:     u.IsAdmin,
		Roles:       []gqlmodel.Role{},
	}
}
