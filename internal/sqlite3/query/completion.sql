-- name: CreateChat :one
INSERT INTO chats (name)
VALUES (?)
RETURNING *;

-- name: CreateMessage :one
INSERT INTO messages (chat_id, content, role_id)
VALUES (?, ?, ?)
RETURNING *;
