package db_helpers

import (
	"context"
	"database/sql"
	"fmt"

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

func GenerateTemplateSliceN(n int64) []*models.Template {
	templates := make([]*models.Template, n, n)
	for i := range n {
		correctedId := i + 1
		templates[i] = &models.Template{
			ID:   correctedId,
			Name: fmt.Sprintf("%d", correctedId),
			Content: fmt.Sprintf(
				`--{{template "%d"}}%d--{{end}}`,
				correctedId,
				correctedId,
			),
		}
	}
	return templates
}
