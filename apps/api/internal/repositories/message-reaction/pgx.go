package message_reaction

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	PgxPool *pgxpool.Pool
}

func NewPgx(opts Opts) *Pgx {
	return &Pgx{
		pgxPool: opts.PgxPool,
	}
}

var _ Repository = &Pgx{}

type Pgx struct {
	pgxPool *pgxpool.Pool
}

var selectFields = []string{
	"r.id",
	"r.message_id",
	"r.user_id",
	"r.reaction",
	"m.channel_id AS channel_id",
}

func (c *Pgx) FindOne(ctx context.Context, reactionID uuid.UUID) (
	*MessageReaction,
	error,
) {
	query, args, err := squirrel.
		Select(selectFields...).
		From("messages_reactions r").
		Where(squirrel.Eq{"r.id": reactionID}).
		Join("chat_messages m ON m.id = r.message_id").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	reaction := &MessageReaction{}

	row := c.pgxPool.QueryRow(ctx, query, args...)
	err = row.Scan(
		&reaction.ID,
		&reaction.MessageID,
		&reaction.UserID,
		&reaction.Reaction,
		&reaction.ChannelID,
	)

	return reaction, err
}

func (c *Pgx) FindManyByMessageID(ctx context.Context, messageID uuid.UUID) (
	[]MessageReaction,
	error,
) {
	query, args, err := squirrel.
		Select(selectFields...).
		From("messages_reactions r").
		Where(squirrel.Eq{"r.message_id": messageID}).
		Join("chat_messages m ON m.id = message_id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pgxPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reactions := make([]MessageReaction, 0)
	for rows.Next() {
		reaction := MessageReaction{}
		err = rows.Scan(
			&reaction.ID,
			&reaction.MessageID,
			&reaction.UserID,
			&reaction.Reaction,
			&reaction.ChannelID,
		)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
	}

	return reactions, nil
}

func (c *Pgx) Create(ctx context.Context, opts CreateMessageReactionOpts) (
	*MessageReaction,
	error,
) {
	query, args, err := squirrel.
		Insert("messages_reactions").
		Columns("message_id", "user_id", "reaction").
		Values(opts.MessageID, opts.UserID, opts.Reaction).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	reaction := &MessageReaction{}
	err = c.pgxPool.QueryRow(ctx, query, args...).Scan(&reaction.ID)
	if err != nil {
		return nil, err
	}

	return c.FindOne(ctx, reaction.ID)
}

func (c *Pgx) Remove(ctx context.Context, id uuid.UUID) error {
	query, args, err := squirrel.
		Delete("messages_reactions").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = c.pgxPool.Exec(ctx, query, args...)
	return err
}
