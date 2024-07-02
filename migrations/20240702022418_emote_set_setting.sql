-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "users" ADD COLUMN "seven_tv_emote_set_id" VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
