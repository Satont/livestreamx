package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateOpts) (*User, error)
	FindByProviderUserID(
		ctx context.Context,
		providerUserID string,
		provider UserConnectionProvider,
	) (
		*User,
		error,
	)
	FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
	Update(ctx context.Context, userID uuid.UUID, opts UpdateOpts) (*User, error)
	AddProviderToUser(ctx context.Context, userID uuid.UUID, opts AddProviderToUserOpts) (
		*User,
		error,
	)
	UpdateProviderByUserID(
		ctx context.Context,
		userID uuid.UUID,
		provider UserConnectionProvider,
		opts UpdateProviderByUserIdOpts,
	) error
	FindByName(ctx context.Context, name string) (*User, error)
}

type CreateOpts struct {
	Name        string
	DisplayName string
	AvatarUrl   string
	Color       string
	Provider    CreateOptsProvider
}

type CreateOptsProvider struct {
	Provider                UserConnectionProvider
	ProviderUserID          string
	ProviderUserName        string
	ProviderUserDisplayName string
	ProviderAvatar          string
}

type UpdateOpts struct {
	Name        *string
	DisplayName *string
	Color       *string
	IsBanned    *bool
}

type AddProviderToUserOpts struct {
	Provider                UserConnectionProvider
	ProviderUserID          string
	ProviderUserName        string
	ProviderUserDisplayName string
	ProviderUserAvatar      string
}

type UpdateProviderByUserIdOpts struct {
	ProviderUserName        string
	ProviderUserDisplayName string
	ProviderAvatar          string
}
