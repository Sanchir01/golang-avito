-- +goose Up
-- +goose StatementBegin
CREATE TYPE city_enum AS ENUM (
    'Москва',
    'Санкт-Петербург',
    'Казань'
    );
CREATE TABLE IF NOT EXISTS pvz(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    registration_date TIMESTAMP NOT NULL DEFAULT NOW(),
    city city_enum NOT NULL,
    version BIGINT NOT NULL DEFAULT 1
);
ALTER TABLE users ADD COLUMN email TEXT NOT NULL UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
