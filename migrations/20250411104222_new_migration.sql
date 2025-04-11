-- +goose Up
-- +goose StatementBegin
ALTER TABLE acceptance DROP CONSTRAINT acceptance_pvz_id_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
