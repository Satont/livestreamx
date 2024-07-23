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
	FindManyByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*User, error)
	FindByStreamKey(ctx context.Context, streamKey uuid.UUID) (*User, error)
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
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
	FindMany(ctx context.Context, opts FindManyOpts) (*FindManyResult, error)
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
	Email                   *string
}

type UpdateOpts struct {
	Name              *string
	DisplayName       *string
	Color             *string
	IsBanned          *bool
	SevenTvEmoteSetID *string
	AvatarUrl         *string
}

type AddProviderToUserOpts struct {
	Provider                UserConnectionProvider
	ProviderUserID          string
	ProviderUserName        string
	ProviderUserDisplayName string
	ProviderUserAvatar      string
	Email                   *string
}

type UpdateProviderByUserIdOpts struct {
	ProviderUserName        string
	ProviderUserDisplayName string
	ProviderAvatar          string
	Email                   *string
}

type FindManyOpts struct {
	Page    int
	PerPage int
}

type FindManyResult struct {
	Users       []User
	Total       int
	HasNextPage bool
}
