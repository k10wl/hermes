package db_helpers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/k10wl/hermes/internal/models"
)

func roleStringToID(role string) int64 {
	switch role {
	case "user":
		return 1
	case "assistant":
		return 2
	case "system":
		return 3
	}
	panic(fmt.Sprintf("role %q does not exist", role))
}

func CreateMessages(db *sql.DB, ctx context.Context, messages []*models.Message) error {
	tx, err := db.BeginTx(ctx, nil)
	vals := []any{}
	if err != nil {
		return err
	}
	sqlStr := "INSERT INTO messages (chat_id, content, role_id) VALUES "
	for _, v := range messages {
		sqlStr += "(?, ?, ?), "
		vals = append(vals, v.ChatID, v.Content, roleStringToID(v.Role))
	}
	sqlStr = strings.TrimSuffix(sqlStr, ", ")
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, vals...)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func GenerateMessagesSliceN(n int64, chatID int64) []*models.Message {
	chats := []*models.Message{}
	for i := 0; i < int(n); i++ {
		id := i + 1
		chats = append(
			chats,
			&models.Message{
				ID:      int64(id),
				ChatID:  chatID,
				Content: "generated",
				Role:    "assistant",
			},
		)
	}
	return chats
}
