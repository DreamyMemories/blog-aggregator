-- +goose up
ALTER TABLE users ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
    encode(digest(random()::text, 'sha256'), 'hex')
);

-- +goose down
ALTER TABLE users DROP COLUMN api_key;