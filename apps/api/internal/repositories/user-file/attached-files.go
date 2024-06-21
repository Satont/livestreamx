package user_file

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateUserFileOpts) (*UserFile, error)
}

type CreateUserFileOpts struct {
	UserID   uuid.UUID
	CdnPath  string
	FileName string
	MimeType string
}

type UserFile struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	CdnPath  string
	FileName string
	MimeType string
}
