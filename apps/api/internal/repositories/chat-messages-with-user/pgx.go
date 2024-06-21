package chat_messages_with_user

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	PgxConn *pgxpool.Pool
}

func NewPgx(opts Opts) *ChatMessageWithUserPgx {
	return &ChatMessageWithUserPgx{
		pgx: opts.PgxConn,
	}
}

var _ Repository = &ChatMessageWithUserPgx{}

type ChatMessageWithUserPgx struct {
	pgx *pgxpool.Pool
}

var selectFields = []string{
	"chat_messages.id",
	"chat_messages.sender_id",
	"chat_messages.text",
	"chat_messages.created_at",
	"users.id",
	"users.name",
	"users.display_name",
	"users.color",
	"users.avatar_url",
	"users.banned",
	"users.created_at",
}

func (c *ChatMessageWithUserPgx) FindByID(ctx context.Context, id uuid.UUID) (
	*MessageWithUser,
	error,
) {
	query, args, err := squirrel.
		Select(selectFields...).
		From("chat_messages chat_messages").
		Join("users users ON chat_messages.sender_id = users.id").
		Where("chat_messages.id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := c.pgx.QueryRow(ctx, query, args...)
	message := &MessageWithUser{}
	err = row.Scan(
		&message.Message.ID,
		&message.Message.SenderID,
		&message.Message.Text,
		&message.Message.CreatedAt,

		&message.User.ID,
		&message.User.Name,
		&message.User.DisplayName,
		&message.User.Color,
		&message.User.AvatarUrl,
		&message.User.Banned,
		&message.User.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (c *ChatMessageWithUserPgx) FindLatest(
	ctx context.Context,
	opts FindLatestOpts,
) ([]MessageWithUser, error) {
	limit := 100
	if opts.Limit != 0 {
		limit = opts.Limit
	}

	if opts.Limit > 500 {
		return nil, errors.New("limit is too high")
	}

	query, args, err := squirrel.
		Select(
			selectFields...,
		).
		From("chat_messages chat_messages").
		Join("users users ON chat_messages.sender_id = users.id").
		OrderBy("chat_messages.created_at ASC").
		Where("users.banned IS FALSE").
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pgx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []MessageWithUser
	for rows.Next() {
		var message MessageWithUser
		err := rows.Scan(
			&message.Message.ID,
			&message.Message.SenderID,
			&message.Message.Text,
			&message.Message.CreatedAt,

			&message.User.ID,
			&message.User.Name,
			&message.User.DisplayName,
			&message.User.Color,
			&message.User.AvatarUrl,
			&message.User.Banned,
			&message.User.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
