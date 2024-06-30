-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE UNIQUE INDEX user_roles_user_id_role_id_unique_index ON user_roles(user_id, role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
