-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    user_email TEXT NOT NULL,
    last_active_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_chats_user_email ON chats(user_email);
CREATE INDEX IF NOT EXISTS idx_chats_last_active_date ON chats(last_active_date);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX IF EXISTS idx_chats_last_active_date;
DROP INDEX IF EXISTS idx_chats_user_email;
DROP TABLE IF EXISTS chats;