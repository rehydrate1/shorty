-- +goose Up
-- SQL section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    short_key VARCHAR(10) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_short_key ON links (short_key);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS links;