package db_helpers

import (
	"context"
	"database/sql"

	"github.com/k10wl/hermes/internal/models"
)

func GetMessagesByChatID(
	db *sql.DB,
	ctx context.Context,
	id int64,
) ([]*models.Message, error) {
	res := []*models.Message{}
	rows, err := db.QueryContext(ctx, `
SELECT
    m.id,
    m.chat_id,
    r.name,
    m.content,
    m.created_at,
    m.updated_at,
    m.deleted_at
FROM
    messages m
JOIN
    roles r ON m.role_id = r.id;
    `)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		msg := models.Message{}
		rows.Scan(
			&msg.ID,
			&msg.ChatID,
			&msg.Role,
			&msg.Content,
			&msg.CreatedAt,
			&msg.UpdatedAt,
			&msg.DeletedAt,
		)
		res = append(res, &msg)
	}
	return res, nil
}
