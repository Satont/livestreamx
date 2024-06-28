package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	DisplayName string     `json:"display_name,omitempty"`
	Color       string     `json:"color,omitempty"`
	AvatarUrl   string     `json:"avatar_url,omitempty"`
	Providers   []Provider `json:"providers,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	Banned      bool       `json:"banned,omitempty"`
	IsAdmin     bool       `json:"is_admin,omitempty"`
	StreamKey   uuid.UUID  `json:"stream_key,omitempty"`
}

type UserConnectionProvider int

const (
	UserConnectionProviderTwitch UserConnectionProvider = iota
	UserConnectionProviderGithub
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
