package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

// UpdateUserProfile is the resolver for the updateUserProfile field.
func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input gqlmodel.UpdateUserProfileInput) (*gqlmodel.User, error) {
	userID, err := r.sessionStorage.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.userRepo.Update(
		ctx,
		uuid.MustParse(userID),
		user.UpdateOpts{
			Color: input.Color.Value(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &gqlmodel.User{
		ID:          user.ID.String(),
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Color:       user.Color,
		Roles:       nil,
		IsBanned:    user.Banned,
		CreatedAt:   user.CreatedAt,
		AvatarURL:   user.AvatarUrl,
	}, nil
}

// UserProfile is the resolver for the userProfile field.
func (r *queryResolver) UserProfile(ctx context.Context) (*gqlmodel.User, error) {
	userID, err := r.sessionStorage.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.userRepo.FindByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}

	return &gqlmodel.User{
		ID:          user.ID.String(),
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Color:       user.Color,
		Roles:       nil,
		IsBanned:    user.Banned,
		CreatedAt:   user.CreatedAt,
		AvatarURL:   user.AvatarUrl,
	}, nil
}
