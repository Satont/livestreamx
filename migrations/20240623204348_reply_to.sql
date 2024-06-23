-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "chat_messages" ADD COLUMN "reply_to" UUID references "chat_messages"("id");
ALTER TABLE "users_providers" ADD COLUMN "email" VARCHAR(500) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
