package db_helpers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/k10wl/hermes/internal/models"
)

type Seeder struct {
	db  *sql.DB
	ctx context.Context
}

func NewSeeder(db *sql.DB, ctx context.Context) *Seeder {
	return &Seeder{
		db:  db,
		ctx: ctx,
	}
}

func (s Seeder) SeedChatsN(n int64) error {
	if n < 0 {
		return fmt.Errorf("cannot process negative N\n")
	}
	return CreateChats(s.db, s.ctx, GenerateChatsSliceN(n))
}

func (s Seeder) SeedMessagesN(n int64, chatID int64) error {
	if n < 0 {
		return fmt.Errorf("cannot process negative N\n")
	}
	return CreateMessages(s.db, s.ctx, GenerateMessagesSliceN(n, chatID))
}

func (s Seeder) SeedTemplatesN(n int64) ([]*models.Template, error) {
	if n < 0 {
		return nil, fmt.Errorf("cannot process negative N\n")
	}

	var err error
	templates := GenerateTemplateSliceN(n)
	for _, template := range templates {
		err = CreateTemplate(s.db, s.ctx, template)
		if err != nil {
			return nil, err
		}
	}
	return templates, nil
}
