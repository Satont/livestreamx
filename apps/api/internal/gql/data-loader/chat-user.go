package data_loader

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

func GetChatUserByID(ctx context.Context, userID uuid.UUID) (*gqlmodel.ChatUser, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.chatUserLoader.Load(ctx, userID)
}

func GetChatUsersByIDs(ctx context.Context, userIDs []uuid.UUID) (
	[]*gqlmodel.ChatUser,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.chatUserLoader.LoadAll(ctx, userIDs)
}

func (c *DataLoader) getChatUsersByIds(ctx context.Context, ids []uuid.UUID) (
	[]*gqlmodel.ChatUser,
	[]error,
) {
	users, err := c.userRepo.FindManyByIDs(ctx, ids)
	if err != nil {
		return nil, []error{err}
	}

	mappedUsers := make([]*gqlmodel.ChatUser, 0, len(users))

	for _, id := range ids {
		foundUser, ok := lo.Find(
			users, func(item *user.User) bool {
				return item.ID == id
			},
		)
		if !ok {
			mappedUsers = append(
				mappedUsers,
				nil,
			)
		} else {
			u := c.mapper.DbUserToChatUser(*foundUser)
			mappedUsers = append(
				mappedUsers,
				&u,
			)
		}
	}

	return mappedUsers, nil
}
