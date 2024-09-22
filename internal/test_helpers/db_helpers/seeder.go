package db_helpers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

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
	chats := []*models.Chat{}
	for i := 0; i < int(n); i++ {
		chats = append(chats, &models.Chat{Name: strconv.Itoa(i)})
	}
	return CreateChats(s.db, s.ctx, chats)
}
