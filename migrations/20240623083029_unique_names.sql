-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE UNIQUE INDEX "users_name_unique_idx" ON "users" ("name");
CREATE UNIQUE INDEX "users_display_name_unique_idx" ON "users" ("display_name");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
