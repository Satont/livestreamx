package user

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestNicknameWithSuffix(t *testing.T) {
	tests := []struct {
		name     string
		nickname string
		try      int
		want     string
	}{
		{name: "first try keeps nickname", nickname: "streamer", try: 0, want: "streamer"},
		{name: "second try adds suffix", nickname: "streamer", try: 1, want: "streamer_2"},
		{name: "third try increments suffix", nickname: "streamer", try: 2, want: "streamer_3"},
		{name: "already suffixed nickname keeps appending", nickname: "streamer_2", try: 1, want: "streamer_2_2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nicknameWithSuffix(tt.nickname, tt.try); got != tt.want {
				t.Fatalf(
					"nicknameWithSuffix(%q, %d) = %q, want %q",
					tt.nickname,
					tt.try,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestUniqueViolationConstraint(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "nil error", err: nil, want: ""},
		{name: "not a pg error", err: errors.New("boom"), want: ""},
		{
			name: "unique violation",
			err:  &pgconn.PgError{Code: "23505", ConstraintName: "users_name_unique_idx"},
			want: "users_name_unique_idx",
		},
		{
			name: "wrapped unique violation",
			err: fmt.Errorf(
				"insert user: %w",
				&pgconn.PgError{Code: "23505", ConstraintName: "users_display_name_unique_idx"},
			),
			want: "users_display_name_unique_idx",
		},
		{
			name: "other pg error",
			err:  &pgconn.PgError{Code: "23503", ConstraintName: "users_id_fkey"},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uniqueViolationConstraint(tt.err); got != tt.want {
				t.Fatalf("uniqueViolationConstraint() = %q, want %q", got, tt.want)
			}
		})
	}
}
