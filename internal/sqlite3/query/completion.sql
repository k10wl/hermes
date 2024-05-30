-- name: CreateChat :one
INSERT INTO chats (name)
VALUES (?)
RETURNING *;

-- name: CreateMessage :one
INSERT INTO messages (chat_id, content, role_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetChats :many
SELECT * FROM chats 
ORDER BY created_at DESC;

-- name: GetChatMessages :many
SELECT
    m.id,
    m.content,
    r.name AS role,  -- Replacing role_id with role name
    m.created_at,
    m.updated_at,
    m.deleted_at
FROM messages m
JOIN roles r ON m.role_id = r.id  -- Join messages with roles based on role_id
WHERE m.chat_id = ?;  -- Replace '?' with the specific chat_id you are interested in

-- name: GetLatestChat :one
SELECT chat_id
FROM messages
ORDER BY created_at DESC
LIMIT 1;
