package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Color       string
	AvatarUrl   string
	Providers   []Provider
	CreatedAt   time.Time
	Banned      bool
	IsAdmin     bool
}

type UserConnectionProvider int

const (
	UserConnectionProviderTwitch UserConnectionProvider = iota
)

type Provider struct {
	ID                      uuid.UUID
	UserID                  uuid.UUID
	Provider                UserConnectionProvider
	ProviderUserID          string
	ProviderUserName        string
	ProviderUserDisplayName string
	ProviderAvatarUrl       string
	Email                   string
}
