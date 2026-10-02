package user

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCreateReassignsTakenNicknames(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("cannot connect to database: %v", err)
	}
	defer pool.Close()

	repo := &Pgx{pgx: pool}

	// Provider ids are unique per test run, so the test does not clash with
	// real data and can clean up only the users it created.
	providerIDPrefix := fmt.Sprintf("nickname-test-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		_, _ = pool.Exec(
			ctx,
			`DELETE FROM users
			WHERE id IN (
				SELECT user_id FROM users_providers WHERE provider_user_id LIKE $1
			)`,
			providerIDPrefix+"%",
		)
	})

	base := fmt.Sprintf("nickname_%d", time.Now().UnixNano())
	email := "user@example.com"
	providerNumber := 0

	create := func(name, displayName string) *User {
		t.Helper()

		providerNumber++
		created, err := repo.Create(ctx, CreateOpts{
			Name:        name,
			DisplayName: displayName,
			AvatarUrl:   "https://example.com/avatar.png",
			Color:       "#ffffff",
			Provider: CreateOptsProvider{
				Provider:                UserConnectionProviderTwitch,
				ProviderUserID:          fmt.Sprintf("%s-%d", providerIDPrefix, providerNumber),
				ProviderUserName:        name,
				ProviderUserDisplayName: displayName,
				ProviderAvatar:          "https://example.com/avatar.png",
				Email:                   &email,
			},
		})
		if err != nil {
			t.Fatalf("cannot create user %q: %v", name, err)
		}

		return created
	}

	assertNickname := func(user *User, name, displayName string) {
		t.Helper()

		if user.Name != name || user.DisplayName != displayName {
			t.Fatalf(
				"got name/display %q/%q, want %q/%q",
				user.Name,
				user.DisplayName,
				name,
				displayName,
			)
		}
	}

	assertNickname(create(base, base), base, base)
	assertNickname(create(base, base), base+"_2", base+"_2")

	// Names are compared case-insensitively, so the next free one is base_3.
	assertNickname(
		create(strings.ToUpper(base), strings.ToUpper(base)),
		strings.ToUpper(base)+"_3",
		strings.ToUpper(base)+"_3",
	)

	// Name is free, but display name is taken, so the pair is shifted together.
	assertNickname(create(base+"x", base), base+"x_4", base+"_4")

	found, err := repo.FindByProviderUserID(
		ctx,
		providerIDPrefix+"-2",
		UserConnectionProviderTwitch,
	)
	if err != nil {
		t.Fatalf("cannot find user by provider id: %v", err)
	}
	assertNickname(found, base+"_2", base+"_2")
}
