-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- drop fk
ALTER TABLE "chat_messages" DROP CONSTRAINT "chat_messages_reply_to_fkey";

ALTER TABLE "chat_messages" ADD CONSTRAINT "fk_chat_messages_reply_to" FOREIGN KEY ("reply_to") REFERENCES "chat_messages"("id") ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
