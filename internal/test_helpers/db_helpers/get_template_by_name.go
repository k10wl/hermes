package db_helpers

import (
	"context"
	"database/sql"

	"github.com/k10wl/hermes/internal/models"
)

func FindTemplateByName(
	db *sql.DB,
	ctx context.Context,
	name string,
) (*models.Template, error) {
	row := db.QueryRowContext(ctx, `
SELECT id, name, content, created_at, updated_at, deleted_at
FROM templates
WHERE name = $1;
    `, name)
	if err := row.Err(); err != nil {
		return nil, err
	}
	template := new(models.Template)
	err := row.Scan(
		&template.ID,
		&template.Name,
		&template.Content,
		&template.Timestamps.CreatedAt,
		&template.Timestamps.UpdatedAt,
		&template.Timestamps.DeletedAt,
	)
	return template, err
}
