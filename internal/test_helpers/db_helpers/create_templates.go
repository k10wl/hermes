package db_helpers

import (
	"context"
	"database/sql"

	"github.com/k10wl/hermes/internal/models"
)

func CreateTemplate(
	db *sql.DB,
	ctx context.Context,
	template *models.Template,
) error {
	_, err := db.Exec(
		"INSERT INTO templates (name, content) VALUES (?, ?)",
		template.Name,
		template.Content,
	)
	return err
}
