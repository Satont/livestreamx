-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- type User struct {
-- 	ID          uuid.UUID
-- 	Name        string
-- 	DisplayName string
-- 	Color       string
-- 	Providers   []Provider
-- 	CreatedAt   time.Time
-- }
--
-- type UserConnectionProvider int
--
-- const (
-- 	UserConnectionProviderTwitch UserConnectionProvider = iota
-- )
--
-- type Provider struct {
-- 	ID                      uuid.UUID
-- 	UserID                  uuid.UUID
-- 	Provider                UserConnectionProvider
-- 	ProviderUserID          string
-- 	ProviderUserName        string
-- 	ProviderUserDisplayName string
-- 	ProviderAvatar          string
-- }

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY default gen_random_uuid(),
    name varchar(255) NOT NULL,
    display_name varchar(255) NOT NULL,
    color varchar(255) NOT NULL,
    created_at TIMESTAMP NOT NULL default now(),
    avatar_url varchar(500) NOT NULL,
    banned BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE users_providers (
    id UUID PRIMARY KEY default gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    provider INT NOT NULL,
    provider_user_id TEXT NOT NULL,
    provider_user_name varchar(255) NOT NULL,
    provider_user_display_name varchar(255) NOT NULL,
    provider_user_avatar_url varchar(500) NOT NULL
);

CREATE TABLE chat_messages (
    id UUID PRIMARY KEY default gen_random_uuid(),
    sender_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL default now()
);

CREATE TABLE channels_roles (
    id UUID PRIMARY KEY default gen_random_uuid(),
    channel_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name varchar(255) NOT NULL,
    image_url varchar(500) NOT NULL,
-- there might be table with prebuilt features and then pivot table to attach those features to role, but i'm to lazy, so that's done via just array of features
    features text[] NOT NULL default '{}'
);

CREATE TABLE user_roles (
    id UUID PRIMARY KEY default gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES channels_roles (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
