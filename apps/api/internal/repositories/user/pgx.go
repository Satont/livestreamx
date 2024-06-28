package user

import (
	"context"
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Pgx *pgxpool.Pool
}

func NewPgx(opts Opts) (*Pgx, error) {
	return &Pgx{
		pgx: opts.Pgx,
	}, nil
}

var _ Repository = &Pgx{}

type Pgx struct {
	pgx *pgxpool.Pool
}

func (c *Pgx) FindByStreamKey(ctx context.Context, streamKey uuid.UUID) (*User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectUserFields...).
		From("users").
		Where(squirrel.Eq{"stream_key": streamKey}).
		Suffix("LIMIT 1").
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = c.pgx.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.DisplayName,
		&user.Color,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.Banned,
		&user.IsAdmin,
		&user.StreamKey,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	providersQuery, providersArgs, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id, user_id, provider, provider_user_id, provider_user_name, provider_user_display_name, provider_user_avatar_url").
		From("users_providers").
		Where(squirrel.Eq{"user_id": user.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pgx.Query(ctx, providersQuery, providersArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []Provider
	for rows.Next() {
		provider := Provider{}
		err = rows.Scan(
			&provider.ID,
			&provider.UserID,
			&provider.Provider,
			&provider.ProviderUserID,
			&provider.ProviderUserName,
			&provider.ProviderUserDisplayName,
			&provider.ProviderAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	user.Providers = providers

	return user, nil
}

func (c *Pgx) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	query, args, err := squirrel.
		Delete("users").
		Where(squirrel.Eq{"id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := c.pgx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

var selectUserFields = []string{
	"id",
	"name",
	"display_name",
	"color",
	"avatar_url",
	"created_at",
	"banned",
	"is_admin",
	"stream_key",
}

func (c *Pgx) FindByName(ctx context.Context, name string) (*User, error) {
	name = strings.ToLower(name)

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectUserFields...).
		From("users").
		Where(
			squirrel.Or{
				squirrel.Expr("LOWER(name) = ?", name),
				squirrel.Expr("LOWER(display_name) = ?", name),
			},
		).
		Suffix("LIMIT 1").
		ToSql()
	if err != nil {
		return nil, err
	}

	row := c.pgx.QueryRow(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	user := &User{}
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.DisplayName,
		&user.Color,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.Banned,
		&user.IsAdmin,
		&user.StreamKey,
	)
	if err != nil {
		return nil, err
	}

	providersQuery, providersArgs, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id, user_id, provider, provider_user_id, provider_user_name, provider_user_display_name, provider_user_avatar_url").
		From("users_providers").
		Where(squirrel.Eq{"user_id": user.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pgx.Query(ctx, providersQuery, providersArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		foundProvider := Provider{}
		err = rows.Scan(
			&foundProvider.ID,
			&foundProvider.UserID,
			&foundProvider.Provider,
			&foundProvider.ProviderUserID,
			&foundProvider.ProviderUserName,
			&foundProvider.ProviderUserDisplayName,
			&foundProvider.ProviderAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		user.Providers = append(user.Providers, foundProvider)
	}

	return user, nil
}

func (c *Pgx) Create(ctx context.Context, opts CreateOpts) (*User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("users").
		Columns(
			"name",
			"display_name",
			"avatar_url",
			"color",
		).
		Values(
			opts.Name,
			opts.DisplayName,
			opts.AvatarUrl,
			opts.Color,
		).
		Suffix("RETURNING id, name, display_name, color, avatar_url, created_at, is_admin, stream_key").
		ToSql()
	if err != nil {
		return nil, err
	}

	tx, err := c.pgx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	user := &User{}
	err = tx.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.DisplayName,
		&user.Color,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.IsAdmin,
		&user.StreamKey,
	)
	if err != nil {
		return nil, err
	}

	providerQuery, providerArgs, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("users_providers").
		Columns(
			"user_id",
			"provider",
			"provider_user_id",
			"provider_user_name",
			"provider_user_display_name",
			"provider_user_avatar_url",
			"email",
		).
		Values(
			user.ID,
			opts.Provider.Provider,
			opts.Provider.ProviderUserID,
			opts.Provider.ProviderUserName,
			opts.Provider.ProviderUserDisplayName,
			opts.Provider.ProviderAvatar,
			opts.Provider.Email,
		).
		Suffix("RETURNING id, user_id, provider, provider_user_id, provider_user_name, provider_user_display_name, provider_user_avatar_url, email").
		ToSql()
	if err != nil {
		return nil, err
	}

	provider := &Provider{}
	err = tx.QueryRow(ctx, providerQuery, providerArgs...).Scan(
		&provider.ID,
		&provider.UserID,
		&provider.Provider,
		&provider.ProviderUserID,
		&provider.ProviderUserName,
		&provider.ProviderUserDisplayName,
		&provider.ProviderAvatarUrl,
		&provider.Email,
	)
	if err != nil {
		return nil, err
	}

	user.Providers = append(user.Providers, *provider)

	return user, nil
}

func (c *Pgx) FindByProviderUserID(
	ctx context.Context,
	providerUserID string,
	provider UserConnectionProvider,
) (*User, error) {
	selectFields := lo.Map(
		selectUserFields, func(field string, _ int) string {
			return "u." + field
		},
	)

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectFields...).
		From("users u").
		Join("users_providers up ON u.id = up.user_id").
		Where(squirrel.Eq{"up.provider_user_id": providerUserID}).
		Where(squirrel.Eq{"up.provider": provider}).
		Suffix("LIMIT 1").
		ToSql()
	if err != nil {
		return nil, err
	}

	row := c.pgx.QueryRow(ctx, query, args...)
	user := &User{}
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.DisplayName,
		&user.Color,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.Banned,
		&user.IsAdmin,
		&user.StreamKey,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	providersQuery, providersArgs, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id, user_id, provider, provider_user_id, provider_user_name, provider_user_display_name, provider_user_avatar_url").
		From("users_providers").
		Where(squirrel.Eq{"user_id": user.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pgx.Query(ctx, providersQuery, providersArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		foundProvider := Provider{}
		err = rows.Scan(
			&foundProvider.ID,
			&foundProvider.UserID,
			&foundProvider.Provider,
			&foundProvider.ProviderUserID,
			&foundProvider.ProviderUserName,
			&foundProvider.ProviderUserDisplayName,
			&foundProvider.ProviderAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		user.Providers = append(user.Providers, foundProvider)
	}

	return user, nil
}

func (c *Pgx) FindByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectUserFields...).
		From("users").
		Where(squirrel.Eq{"id": userID}).
		Suffix("LIMIT 1").
		ToSql()
	if err != nil {
		return nil, err
	}

	row := c.pgx.QueryRow(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	user := &User{}
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.DisplayName,
		&user.Color,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.Banned,
		&user.IsAdmin,
		&user.StreamKey,
	)
	if err != nil {
		return nil, err
	}

	providersQuery, providersArgs, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id, user_id, provider, provider_user_id, provider_user_name, provider_user_display_name, provider_user_avatar_url").
		From("users_providers").
		Where(squirrel.Eq{"user_id": user.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pgx.Query(ctx, providersQuery, providersArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		foundProvider := Provider{}
		err = rows.Scan(
			&foundProvider.ID,
			&foundProvider.UserID,
			&foundProvider.Provider,
			&foundProvider.ProviderUserID,
			&foundProvider.ProviderUserName,
			&foundProvider.ProviderUserDisplayName,
			&foundProvider.ProviderAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		user.Providers = append(user.Providers, foundProvider)
	}

	return user, nil
}

func (c *Pgx) Update(ctx context.Context, userID uuid.UUID, opts UpdateOpts) (*User, error) {
	var updateMap = map[string]interface{}{}

	if opts.Name != nil {
		updateMap["name"] = *opts.Name
	}

	if opts.DisplayName != nil {
		updateMap["display_name"] = *opts.DisplayName
	}

	if opts.Color != nil {
		updateMap["color"] = *opts.Color
	}

	if opts.IsBanned != nil {
		updateMap["banned"] = *opts.IsBanned
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("users").
		SetMap(updateMap).
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	_, err = c.pgx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return c.FindByID(ctx, userID)
}

func (c *Pgx) AddProviderToUser(
	ctx context.Context,
	userID uuid.UUID,
	opts AddProviderToUserOpts,
) (*User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("users_providers").
		Columns(
			"user_id",
			"provider",
			"provider_user_id",
			"provider_user_name",
			"provider_user_display_name",
			"provider_user_avatar_url",
			"email",
		).
		Values(
			userID,
			opts.Provider,
			opts.ProviderUserID,
			opts.ProviderUserName,
			opts.ProviderUserDisplayName,
			opts.ProviderUserAvatar,
			opts.Email,
		).
		ToSql()
	if err != nil {
		return nil, err
	}

	_, err = c.pgx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return c.FindByID(ctx, userID)
}

func (c *Pgx) UpdateProviderByUserID(
	ctx context.Context,
	userID uuid.UUID,
	provider UserConnectionProvider,
	opts UpdateProviderByUserIdOpts,
) error {
	updateMap := map[string]interface{}{
		"provider_user_name":         opts.ProviderUserName,
		"provider_user_display_name": opts.ProviderUserDisplayName,
		"provider_user_avatar_url":   opts.ProviderAvatar,
		"email":                      opts.Email,
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("users_providers").
		SetMap(updateMap).
		Where(
			squirrel.Eq{
				"user_id":  userID,
				"provider": provider,
			},
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = c.pgx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
