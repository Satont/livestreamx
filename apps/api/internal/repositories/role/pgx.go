package role

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	PgxPool *pgxpool.Pool
}

func NewPgx(opts Opts) *RolePgx {
	return &RolePgx{
		pgxPool: opts.PgxPool,
	}
}

var _ Repository = &RolePgx{}

type RolePgx struct {
	pgxPool *pgxpool.Pool
}

func (c *RolePgx) FindUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	query, args, err := squirrel.
		Select("r.id", "r.channel_id", "r.name", "r.image_url", "r.features").
		From("channels_roles r").
		Join("user_roles ur ON r.id = ur.role_id").
		Where(squirrel.Eq{"ur.user_id": userID}).
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

	var roles []Role
	for rows.Next() {
		role := Role{}
		err := rows.Scan(
			&role.ID,
			&role.ChannelID,
			&role.Name,
			&role.ImageUrl,
			&role.Features,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (c *RolePgx) AssignRoleToUser(ctx context.Context, roleID, userID uuid.UUID) error {
	query, args, err := squirrel.
		Insert("user_roles").
		Columns("role_id", "user_id").
		Values(roleID, userID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = c.pgxPool.Exec(ctx, query, args...)
	return err
}

func (c *RolePgx) UnassignRoleFromUser(ctx context.Context, roleID, userID uuid.UUID) error {
	query, args, err := squirrel.
		Delete("user_roles").
		Where(squirrel.Eq{"role_id": roleID}).
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = c.pgxPool.Exec(ctx, query, args...)
	return err
}

func (c *RolePgx) UpdateByID(ctx context.Context, id uuid.UUID, opts UpdateOpts) (*Role, error) {
	query := squirrel.Update("channels_roles").Where(squirrel.Eq{"id": id})

	if opts.Name != nil {
		query = query.Set("name", *opts.Name)
	}
	if opts.ImageUrl != nil {
		query = query.Set("image_url", *opts.ImageUrl)
	}
	if opts.Features != nil {
		query = query.Set("features", opts.Features)
	}

	queryStr, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	if _, err := c.pgxPool.Exec(ctx, queryStr, args...); err != nil {
		return nil, err
	}

	return c.FindOneByID(ctx, id)
}

var selectFields = []string{
	"id",
	"channel_id",
	"name",
	"image_url",
	"features",
}

func (c *RolePgx) Create(ctx context.Context, opts CreateOpts) (*Role, error) {
	query, args, err := squirrel.
		Insert("channels_roles").
		Columns("channel_id", "name", "image_url", "features").
		Values(opts.ChannelID, opts.Name, opts.ImageUrl, opts.Features).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	if err := c.pgxPool.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return nil, err
	}

	return c.FindOneByID(ctx, id)
}

func (c *RolePgx) FindManyByChannelID(ctx context.Context, channelID uuid.UUID) ([]Role, error) {
	query, args, err := squirrel.
		Select(selectFields...).
		From("channels_roles").
		Where(squirrel.Eq{"channel_id": channelID}).
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

	var roles []Role
	for rows.Next() {
		role := Role{}
		err := rows.Scan(
			&role.ID,
			&role.ChannelID,
			&role.Name,
			&role.ImageUrl,
			&role.Features,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (c *RolePgx) FindOneByID(ctx context.Context, id uuid.UUID) (*Role, error) {
	query, args, err := squirrel.
		Select(selectFields...).
		From("channels_roles").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	role := &Role{}

	row := c.pgxPool.QueryRow(ctx, query, args...)
	err = row.Scan(
		&role.ID,
		&role.ChannelID,
		&role.Name,
		&role.ImageUrl,
		&role.Features,
	)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (c *RolePgx) DeleteByID(ctx context.Context, id uuid.UUID) error {
	cmd, err := c.pgxPool.Exec(ctx, "DELETE FROM channels_roles WHERE id = $1", id)
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("role with id %s not found", id)
	}

	return err
}
