-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "messages_reactions" (
    "id" UUID PRIMARY KEY default gen_random_uuid(),
    "message_id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "reaction" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY ("message_id") REFERENCES "chat_messages" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

CREATE INDEX "messages_reactions_message_id_index" ON "messages_reactions" ("message_id");
CREATE UNIQUE INDEX "messages_reactions_message_id_user_id_reaction_index" ON "messages_reactions" ("message_id", "user_id", "reaction");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
