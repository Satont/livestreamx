package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
	userrepo "github.com/satont/stream/apps/api/internal/repositories/user"
)

type callbackOpts struct {
	Provider                userrepo.UserConnectionProvider
	ProviderUserID          string
	ProviderUserName        string
	ProviderUserDisplayName string
	ProviderAvatar          string
	Email                   *string
}

// case 1: user not logged in and user with provider user id doesnt exists => create user
// case 2: user not logged in and user with provider user id exists => login user
// case 3: user logged in and user with provider user id doesnt exists => link user
// case 4: user logged in and provider linked to another user => throw error

func (c *Auth) callback(ctx context.Context, opts callbackOpts) (*userrepo.User, error) {
	authedUserID, authedErr := c.sessionStore.GetUserID(ctx)
	existedUser, _ := c.userRepo.FindByProviderUserID(
		ctx,
		opts.ProviderUserID,
		opts.Provider,
	)

	subOpts := subCallbackOpts{
		callbackOpts:          opts,
		existedUserByProvider: existedUser,
		authedUserID:          authedUserID,
	}

	if authedErr != nil && existedUser == nil { // user not logged and user by provider doesnt exists => create user
		return c.callbackHandleCreate(ctx, subOpts)
	} else if authedErr != nil && existedUser != nil { // user not logged in and user by provider exists => login
		return existedUser, nil
	} else if authedErr == nil && existedUser == nil { // user logged in and user by provider doesnt exists => link provider to user
		return c.callbackHandleLink(ctx, subOpts)
	} else if existedUser != nil && authedUserID != existedUser.ID.String() { // user logged in and provider linked to another user => throw error
		return nil, fmt.Errorf("user with provider user id %s already exists", opts.ProviderUserID)
	} else {
		return nil, fmt.Errorf("unexpected error")
	}

	return nil, fmt.Errorf("unexpected error")
}

type subCallbackOpts struct {
	callbackOpts

	authedUserID          string
	existedUserByProvider *userrepo.User
}

func (c *Auth) callbackHandleLink(ctx context.Context, opts subCallbackOpts) (
	*userrepo.User,
	error,
) {
	return c.userRepo.AddProviderToUser(
		ctx,
		uuid.MustParse(opts.authedUserID),
		userrepo.AddProviderToUserOpts{
			Provider:                opts.callbackOpts.Provider,
			ProviderUserID:          opts.callbackOpts.ProviderUserID,
			ProviderUserName:        opts.callbackOpts.ProviderUserName,
			ProviderUserDisplayName: opts.callbackOpts.ProviderUserDisplayName,
			ProviderUserAvatar:      opts.callbackOpts.ProviderAvatar,
			Email:                   opts.callbackOpts.Email,
		},
	)
}

func (c *Auth) callbackHandleCreate(ctx context.Context, opts subCallbackOpts) (
	*userrepo.User,
	error,
) {
	return c.userRepo.Create(
		ctx,
		userrepo.CreateOpts{
			Name:        opts.ProviderUserName,
			DisplayName: opts.ProviderUserDisplayName,
			AvatarUrl:   opts.ProviderAvatar,
			Color:       colorful.WarmColor().Hex(),
			Provider: userrepo.CreateOptsProvider{
				Provider:                opts.Provider,
				ProviderUserID:          opts.ProviderUserID,
				ProviderUserName:        opts.ProviderUserName,
				ProviderUserDisplayName: opts.ProviderUserDisplayName,
				ProviderAvatar:          opts.ProviderAvatar,
				Email:                   opts.Email,
			},
		},
	)
}
