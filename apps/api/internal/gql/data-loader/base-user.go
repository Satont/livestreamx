package data_loader

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

func GetBaseUserByID(ctx context.Context, userID uuid.UUID) (*gqlmodel.BaseUser, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.baseUserLoader.Load(ctx, userID)
}

func GetBaseUsersByIDs(ctx context.Context, userIDs []uuid.UUID) (
	[]*gqlmodel.BaseUser,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.baseUserLoader.LoadAll(ctx, userIDs)
}

func (c *DataLoader) getBaseUsersByIds(ctx context.Context, ids []uuid.UUID) (
	[]*gqlmodel.BaseUser,
	[]error,
) {
	users, err := c.userRepo.FindManyByIDs(ctx, ids)
	if err != nil {
		return nil, []error{err}
	}

	mappedUsers := make([]*gqlmodel.BaseUser, 0, len(users))

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
			u := c.mapper.DbUserToBaseUserGql(*foundUser)
			mappedUsers = append(
				mappedUsers,
				&u,
			)
		}
	}

	return mappedUsers, nil
}
