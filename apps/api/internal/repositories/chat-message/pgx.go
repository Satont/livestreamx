package chat_message

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (c *ChatMessagePgx) Create(ctx context.Context, opts CreateChatMessageOpts) (
	*Message,
	error,
) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("chat_messages").
		Columns("id", "sender_id", "text", "reply_to", "channel_id").
		Values(uuid.New(), opts.SenderID, opts.Text, opts.ReplyTo, opts.ChannelID).
		Suffix("RETURNING id, sender_id, text, created_at, reply_to, channel_id").
		ToSql()
	if err != nil {
		return nil, err
	}

	message := Message{}
	if err := c.pgx.QueryRow(ctx, query, args...).Scan(
		&message.ID,
		&message.SenderID,
		&message.Text,
		&message.CreatedAt,
		&message.ReplyTo,
		&message.ChannelID,
	); err != nil {
		return nil, err
	}

	return &message, nil
}

func (c *ChatMessagePgx) FindLatest(ctx context.Context, opts FindManyOpts) ([]Message, error) {
	limit := 100
	if opts.Limit > 0 {
		limit = opts.Limit
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "sender_id", "text", "created_at", "channel_id").
		From("chat_messages").
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
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
		if err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.Text,
			&message.CreatedAt,
			&message.ChannelID,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
