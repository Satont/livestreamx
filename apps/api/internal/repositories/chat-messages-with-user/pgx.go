package chat_messages_with_user

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	message_reaction "github.com/satont/stream/apps/api/internal/repositories/message-reaction"
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
	"chat_messages.reply_to",
	"users.id",
	"users.name",
	"users.display_name",
	"users.color",
	"users.avatar_url",
	"users.banned",
	"users.created_at",
	"users.is_admin",
	"mr.id AS reaction_id",
	"mr.message_id AS reaction_message_id",
	"mr.user_id AS reaction_user_id",
	"mr.reaction",
	"mr.created_at AS reaction_created_at",
}

type tempMessageReaction struct {
	ID        *uuid.UUID
	MessageID *uuid.UUID
	UserID    *uuid.UUID
	Reaction  *string
	CreatedAt *time.Time
}

func (c *ChatMessageWithUserPgx) FindByID(ctx context.Context, id uuid.UUID) (
	*MessageWithUser,
	error,
) {
	query, args, err := squirrel.
		Select(selectFields...).
		From("chat_messages chat_messages").
		Join("users users ON chat_messages.sender_id = users.id").
		LeftJoin("messages_reactions mr ON chat_messages.id = mr.message_id").
		Where("chat_messages.id = ?", id).
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

	var message *MessageWithUser
	for rows.Next() {
		row := &MessageWithUser{}
		reaction := &tempMessageReaction{}

		err = rows.Scan(
			&row.Message.ID,
			&row.Message.SenderID,
			&row.Message.Text,
			&row.Message.CreatedAt,
			&row.Message.ReplyTo,

			&row.User.ID,
			&row.User.Name,
			&row.User.DisplayName,
			&row.User.Color,
			&row.User.AvatarUrl,
			&row.User.Banned,
			&row.User.CreatedAt,
			&row.User.IsAdmin,

			&reaction.ID,
			&reaction.MessageID,
			&reaction.UserID,
			&reaction.Reaction,
			&reaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if message == nil {
			message = row
		}

		if reaction.ID != nil {
			message.Reactions = append(
				message.Reactions,
				message_reaction.MessageReaction{
					ID:        *reaction.ID,
					MessageID: *reaction.MessageID,
					UserID:    *reaction.UserID,
					Reaction:  *reaction.Reaction,
					CreatedAt: *reaction.CreatedAt,
				},
			)
		}
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
		OrderBy("chat_messages.created_at DESC").
		Where("users.banned IS FALSE").
		LeftJoin("messages_reactions mr ON chat_messages.id = mr.message_id").
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
		message := MessageWithUser{}
		reaction := &tempMessageReaction{}

		err = rows.Scan(
			&message.Message.ID,
			&message.Message.SenderID,
			&message.Message.Text,
			&message.Message.CreatedAt,
			&message.Message.ReplyTo,

			&message.User.ID,
			&message.User.Name,
			&message.User.DisplayName,
			&message.User.Color,
			&message.User.AvatarUrl,
			&message.User.Banned,
			&message.User.CreatedAt,
			&message.User.IsAdmin,

			&reaction.ID,
			&reaction.MessageID,
			&reaction.UserID,
			&reaction.Reaction,
			&reaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		foundMessageIndex := -1
		for i, m := range messages {
			if m.Message.ID == message.Message.ID {
				foundMessageIndex = i
				break
			}
		}

		if foundMessageIndex == -1 {
			if reaction.ID != nil {
				message.Reactions = append(
					message.Reactions,
					message_reaction.MessageReaction{
						ID:        *reaction.ID,
						MessageID: *reaction.MessageID,
						UserID:    *reaction.UserID,
						Reaction:  *reaction.Reaction,
						CreatedAt: *reaction.CreatedAt,
					},
				)
			}
			messages = append(messages, message)
		} else {
			if reaction.ID != nil {
				messages[foundMessageIndex].Reactions = append(
					messages[foundMessageIndex].Reactions,
					message_reaction.MessageReaction{
						ID:        *reaction.ID,
						MessageID: *reaction.MessageID,
						UserID:    *reaction.UserID,
						Reaction:  *reaction.Reaction,
						CreatedAt: *reaction.CreatedAt,
					},
				)
			}
		}
	}

	return messages, nil
}
