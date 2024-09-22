package db_helpers

import (
	"context"
	"database/sql"
	"strings"

	"github.com/k10wl/hermes/internal/models"
)

func CreateChats(db *sql.DB, ctx context.Context, chats []*models.Chat) error {
	tx, err := db.BeginTx(ctx, nil)
	vals := []any{}
	if err != nil {
		return err
	}
	sqlStr := "INSERT INTO chats (name) VALUES "
	for _, v := range chats {
		sqlStr += "(?), "
		vals = append(vals, v.Name)
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
