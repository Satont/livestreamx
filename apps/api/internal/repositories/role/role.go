package role

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateOpts) (*Role, error)
	FindMany(ctx context.Context) ([]Role, error)
	FindOneByID(ctx context.Context, id uuid.UUID) (*Role, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	UpdateByID(ctx context.Context, id uuid.UUID, opts UpdateOpts) (*Role, error)
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
