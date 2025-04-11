-- +goose Up
-- +goose StatementBegin
CREATE TYPE acceptance_status AS ENUM (
    'close',
    'in_progress'
);

CREATE TABLE IF NOT EXISTS acceptance(
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  pvz_id UUID NOT NULL UNIQUE,
  version BIGINT NOT NULL DEFAULT 1,

  CONSTRAINT fk_pvz_receiving FOREIGN KEY (pvz_id) REFERENCES pvz(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product(
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  version BIGINT NOT NULL DEFAULT 1,
  receiving_id UUID NOT NULL,

  CONSTRAINT fk_item_receiving FOREIGN KEY (receiving_id) REFERENCES acceptance(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
