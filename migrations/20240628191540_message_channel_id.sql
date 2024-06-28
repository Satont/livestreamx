-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE chat_messages ADD COLUMN channel_id UUID NOT NULL references users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
