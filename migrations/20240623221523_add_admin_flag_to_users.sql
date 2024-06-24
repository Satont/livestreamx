-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "users" ADD COLUMN "is_admin" BOOLEAN NOT NULL default false;
alter table channels_roles
    alter column image_url type varchar(1000) using image_url::varchar(1000);

alter table channels_roles
    alter column image_url drop not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
