-- +goose Up
-- +goose StatementBegin
CREATE TYPE product_type_enum AS ENUM (
    'электроника',
    'одежда',
    'обувь'
);
ALTER TABLE product ADD COLUMN type product_type_enum NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
