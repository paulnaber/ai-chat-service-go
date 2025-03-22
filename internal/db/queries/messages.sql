-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 LIMIT 1;

-- name: GetMessagesByChatID :many
SELECT * FROM messages
WHERE chat_id = $1
ORDER BY created_at ASC;

-- name: CreateMessage :one
INSERT INTO messages (id, content, sender_type, chat_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;