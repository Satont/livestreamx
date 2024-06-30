package role

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateOpts) (*Role, error)
	FindManyByChannelID(ctx context.Context, channelID uuid.UUID) ([]Role, error)
	FindOneByID(ctx context.Context, id uuid.UUID) (*Role, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	UpdateByID(ctx context.Context, id uuid.UUID, opts UpdateOpts) (*Role, error)

	FindUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error)
	AssignRoleToUser(ctx context.Context, roleID, userID uuid.UUID) error
	UnassignRoleFromUser(ctx context.Context, roleID, userID uuid.UUID) error
}

type CreateOpts struct {
	ChannelID uuid.UUID
	Name      string
	ImageUrl  *string
	Features  []string
}

type UpdateOpts struct {
	Name     *string
	ImageUrl *string
	Features *[]string
}
