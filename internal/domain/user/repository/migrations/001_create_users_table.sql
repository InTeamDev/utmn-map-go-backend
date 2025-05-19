-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    modified_at TIMESTAMP WITH TIME ZONE NOT NULL,
    photo_url TEXT
);

-- Create an index on telegram_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);

-- Create an index on username for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- +migrate Down
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_telegram_id;
DROP TABLE IF EXISTS users; 