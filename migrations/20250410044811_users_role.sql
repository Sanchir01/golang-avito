-- +goose Up
-- +goose StatementBegin
CREATE TYPE role_enum AS ENUM (
    'employee',
    'moderator'
);
ALTER TABLE users ADD COLUMN role role_enum NOT NULL DEFAULT 'employee';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
