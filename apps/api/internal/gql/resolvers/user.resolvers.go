package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	req "github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/httpserver/middlewares"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

// UpdateUserProfile is the resolver for the updateUserProfile field.
func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input gqlmodel.UpdateUserProfileInput) (*gqlmodel.AuthedUser, error) {
	currentUser := middlewares.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	opts := user.UpdateOpts{
		SevenTvEmoteSetID: input.SevenTvEmoteSetID.Value(),
	}
	if input.Color.IsSet() {
		opts.Color = input.Color.Value()
	}

	if input.AvatarURL.IsSet() {
		resp, err := req.Get(*input.AvatarURL.Value())
		if err != nil {
			return nil, fmt.Errorf("cannot fetch avatar url: %w", err)
		}

		if !strings.HasPrefix(resp.GetContentType(), "image") {
			return nil, fmt.Errorf("provided avatar is not an image")
		}

		opts.AvatarUrl = input.AvatarURL.Value()
	}

	if input.Name.IsSet() && input.DisplayName.IsSet() {
		if strings.ToLower(*input.Name.Value()) != strings.ToLower(*input.DisplayName.Value()) {
			return nil, fmt.Errorf("name and display does not match")
		}

		if !userNameRegexp.MatchString(*input.Name.Value()) {
			return nil, fmt.Errorf("name does not match pattern: %s", userNameRegexp.String())
		}

		if utf8.RuneCountInString(*input.Name.Value()) < 3 || utf8.RuneCountInString(*input.Name.Value()) > 25 {
			return nil, fmt.Errorf("name length must be between 3 and 25 characters")
		}

		if utf8.RuneCountInString(*input.DisplayName.Value()) < 3 || utf8.RuneCountInString(*input.DisplayName.Value()) > 25 {
			return nil, fmt.Errorf("display name length must be between 3 and 25 characters")
		}

		if !userNameRegexp.MatchString(*input.DisplayName.Value()) {
			return nil, fmt.Errorf(
				"display name does not match pattern: %s",
				userNameRegexp.String(),
			)
		}

		opts.Name = lo.ToPtr(strings.ToLower(*input.Name.Value()))
		opts.DisplayName = input.DisplayName.Value()
	}

	newUser, err := r.userRepo.Update(
		ctx,
		currentUser.ID,
		opts,
	)
	if err != nil {
		return nil, err
	}

	providers := make([]gqlmodel.AuthedUserProvider, 0, len(newUser.Providers))
	for _, provider := range newUser.Providers {
		providers = append(
			providers, gqlmodel.AuthedUserProvider{
				Provider:    r.mapper.DbProviderToGql(provider.Provider),
				UserID:      provider.ProviderUserID,
				Name:        provider.ProviderUserName,
				DisplayName: provider.ProviderUserDisplayName,
				AvatarURL:   provider.ProviderAvatarUrl,
			},
		)
	}

	go func() {
		if err := r.sevenTv.InitUser(*newUser); err != nil {
			r.logger.Sugar().Error("Cannot init 7tv user", err)
		}
	}()

	return &gqlmodel.AuthedUser{
		ID:                newUser.ID,
		Name:              newUser.Name,
		DisplayName:       newUser.DisplayName,
		Color:             newUser.Color,
		IsBanned:          newUser.Banned,
		CreatedAt:         newUser.CreatedAt,
		AvatarURL:         newUser.AvatarUrl,
		IsAdmin:           newUser.IsAdmin,
		Providers:         providers,
		StreamKey:         newUser.StreamKey,
		SevenTvEmoteSetID: newUser.SevenTvEmoteSetID,
	}, nil
}

// DeleteAccount is the resolver for the deleteAccount field.
func (r *mutationResolver) DeleteAccount(ctx context.Context) (bool, error) {
	userID, err := r.sessionStorage.GetUserID(ctx)
	if err != nil {
		return false, err
	}

	if _, err = r.Logout(ctx); err != nil {
		return false, err
	}

	err = r.userRepo.DeleteAccount(ctx, uuid.MustParse(userID))
	if err != nil {
		return false, err
	}

	return true, nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	session := r.sessionStorage.GetSession(ctx)
	if session == nil {
		return false, nil
	}

	session.Clear()
	session.Save()

	return true, nil
}

// UserProfile is the resolver for the userProfile field.
func (r *queryResolver) UserProfile(ctx context.Context) (*gqlmodel.AuthedUser, error) {
	userID, err := r.sessionStorage.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.userRepo.FindByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}

	providers := make([]gqlmodel.AuthedUserProvider, 0, len(user.Providers))
	for _, provider := range user.Providers {
		providers = append(
			providers, gqlmodel.AuthedUserProvider{
				Provider:    r.mapper.DbProviderToGql(provider.Provider),
				UserID:      provider.ProviderUserID,
				Name:        provider.ProviderUserName,
				DisplayName: provider.ProviderUserDisplayName,
				AvatarURL:   provider.ProviderAvatarUrl,
			},
		)
	}

	return &gqlmodel.AuthedUser{
		ID:                user.ID,
		Name:              user.Name,
		DisplayName:       user.DisplayName,
		Color:             user.Color,
		IsBanned:          user.Banned,
		CreatedAt:         user.CreatedAt,
		AvatarURL:         user.AvatarUrl,
		IsAdmin:           user.IsAdmin,
		Providers:         providers,
		StreamKey:         user.StreamKey,
		SevenTvEmoteSetID: user.SevenTvEmoteSetID,
	}, nil
}

// FetchUserByName is the resolver for the fetchUserByName field.
func (r *queryResolver) FetchUserByName(ctx context.Context, name string) (gqlmodel.User, error) {
	user, err := r.userRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return r.mapper.DbUserToBaseUserGql(*user), nil
}

// FetchUserByID is the resolver for the fetchUserById field.
func (r *queryResolver) FetchUserByID(ctx context.Context, id uuid.UUID) (gqlmodel.User, error) {
	user, err := r.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapper.DbUserToBaseUserGql(*user), nil
}
