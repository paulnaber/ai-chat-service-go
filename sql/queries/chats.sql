-- name: GetChat :one
SELECT * FROM chats
WHERE id = $1 LIMIT 1;

-- name: GetChatsByUserEmail :many
SELECT * FROM chats
WHERE user_email = $1
ORDER BY last_active_date DESC;

-- name: CreateChat :one
INSERT INTO chats (id, title, user_email, last_active_date, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateChatLastActive :exec
UPDATE chats
SET last_active_date = $2, updated_at = $3
WHERE id = $1;
