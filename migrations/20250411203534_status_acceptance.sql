-- +goose Up
-- +goose StatementBegin
ALTER TABLE acceptance ADD COLUMN status TEXT DEFAULT 'in_progress' NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
