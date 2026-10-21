-- +goose Up
-- +goose StatementBegin
ALTER TABLE chat_messages ADD COLUMN sender_color VARCHAR(255);

UPDATE chat_messages cm
SET sender_color = u.color
FROM users u
WHERE u.id = cm.sender_id;

ALTER TABLE chat_messages ALTER COLUMN sender_color SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chat_messages DROP COLUMN sender_color;
-- +goose StatementEnd
