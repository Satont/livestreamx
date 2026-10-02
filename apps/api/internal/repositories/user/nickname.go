package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

const maxNicknameTries = 10

// nicknameWithSuffix returns the nickname itself for the first try and appends
// a numeric suffix for the next ones: "name" -> "name_2", "name" -> "name_3".
func nicknameWithSuffix(nickname string, try int) string {
	if try <= 0 {
		return nickname
	}

	return fmt.Sprintf("%s_%d", nickname, try+1)
}

// pickFreeNickname finds the first name/display name pair which is not taken by
// another user. Names are compared case-insensitively, the same way they are
// resolved by FindByName, so mentions never point to two users at once.
func (c *Pgx) pickFreeNickname(
	ctx context.Context,
	name string,
	displayName string,
	startTry int,
) (string, string, error) {
	for try := startTry; try < startTry+maxNicknameTries; try++ {
		candidateName := nicknameWithSuffix(name, try)
		candidateDisplayName := nicknameWithSuffix(displayName, try)

		taken, err := c.nicknamesTaken(ctx, candidateName, candidateDisplayName)
		if err != nil {
			return "", "", err
		}

		if !taken {
			return candidateName, candidateDisplayName, nil
		}
	}

	return "", "", fmt.Errorf("could not find a free nickname for %q", name)
}

func (c *Pgx) nicknamesTaken(ctx context.Context, name, displayName string) (bool, error) {
	const query = `
		SELECT
			EXISTS (
				SELECT 1
				FROM users
				WHERE LOWER(name) = LOWER($1) OR LOWER(display_name) = LOWER($1)
			) OR EXISTS (
				SELECT 1
				FROM users
				WHERE LOWER(name) = LOWER($2) OR LOWER(display_name) = LOWER($2)
			)
	`

	var taken bool
	if err := c.pgx.QueryRow(ctx, query, name, displayName).Scan(&taken); err != nil {
		return false, err
	}

	return taken, nil
}

// uniqueViolationConstraint returns the name of the violated unique constraint
// or an empty string if the error is not a unique violation.
func uniqueViolationConstraint(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr.ConstraintName
	}

	return ""
}
