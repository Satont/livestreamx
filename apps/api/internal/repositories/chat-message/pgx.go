package chat_message

import (
	"context"
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

func NewPgx(opts Opts) *ChatMessagePgx {
	return &ChatMessagePgx{
		pgx: opts.PgxConn,
	}
}

var _ Repository = &ChatMessagePgx{}

type ChatMessagePgx struct {
	pgx *pgxpool.Pool
}

type tempMessageReaction struct {
	ID        *uuid.UUID
	MessageID *uuid.UUID
	UserID    *uuid.UUID
	Reaction  *string
	CreatedAt *time.Time
}

var selectColumns = []string{
	"cm.id",
	"cm.sender_id",
	"cm.text",
	"cm.created_at",
	"cm.reply_to",
	"cm.channel_id",
	"mr.id AS reaction_id",
	"mr.message_id AS reaction_message_id",
	"mr.user_id AS reaction_user_id",
	"mr.reaction",
	"mr.created_at AS reaction_created_at",
}

func (c *ChatMessagePgx) FindByID(ctx context.Context, id uuid.UUID) (*Message, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectColumns...).
		From("chat_messages cm").
		LeftJoin("messages_reactions mr ON cm.id = mr.message_id").
		Where(squirrel.Eq{"cm.id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pgx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var message *Message
	for rows.Next() {
		row := &Message{}
		reaction := &tempMessageReaction{}

		err = rows.Scan(
			&row.ID,
			&row.SenderID,
			&row.Text,
			&row.CreatedAt,
			&row.ReplyTo,
			&row.ChannelID,

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
					ChannelID: row.ChannelID,
				},
			)
		}
	}

	return message, nil
}

func (c *ChatMessagePgx) Create(ctx context.Context, opts CreateChatMessageOpts) (*Message, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("chat_messages").
		Columns("sender_id", "text", "reply_to", "channel_id").
		Values(opts.SenderID, opts.Text, opts.ReplyTo, opts.ChannelID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	err = c.pgx.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, err
	}

	return c.FindByID(ctx, id)
}

func (c *ChatMessagePgx) FindLatest(ctx context.Context, opts FindLatestOpts) ([]Message, error) {
	limit := 100
	if opts.Limit > 0 {
		limit = opts.Limit
	}

	query, args, err := squirrel.
		Select(selectColumns...).
		From("chat_messages AS cm").
		LeftJoin("messages_reactions mr ON cm.id = mr.message_id").
		OrderBy("created_at DESC").
		Where(
			squirrel.Eq{
				"cm.channel_id": opts.ChannelID,
			},
		).
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

	var messages []Message
	for rows.Next() {
		message := Message{}
		reaction := &tempMessageReaction{}

		err = rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.Text,
			&message.CreatedAt,
			&message.ReplyTo,
			&message.ChannelID,

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
			if m.ID == message.ID {
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
						ChannelID: message.ChannelID,
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
						ChannelID: message.ChannelID,
					},
				)
			}
		}
	}

	return messages, nil
}
