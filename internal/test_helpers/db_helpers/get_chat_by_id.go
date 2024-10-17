package db_helpers

import (
	"context"
	"database/sql"

	"github.com/k10wl/hermes/internal/models"
)

func GetChatByID(db *sql.DB, ctx context.Context, id int64) (*models.Chat, error) {
	var chat models.Chat
	row := db.QueryRowContext(ctx, `
SELECT 
    id,
    name,
    created_at,
    updated_at,
    deleted_at
FROM chats
WHERE id = ?`, id)
	if err := row.Err(); err != nil {
		return nil, err
	}
	err := row.Scan(
		&chat.ID,
		&chat.Name,
		&chat.CreatedAt,
		&chat.UpdatedAt,
		&chat.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}
