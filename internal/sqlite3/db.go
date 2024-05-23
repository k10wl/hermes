package sqlite3

import (
	"context"
	"database/sql"
	"path"

	"github.com/k10wl/hermes/internal/runtime"
	"github.com/k10wl/hermes/internal/sqlc"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite3 struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewSQLite3(config *runtime.Config) (*SQLite3, error) {
	dbName := path.Join(config.ConfigDir, "main.db")
	err := runMigrations(dbName)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(db)
	return &SQLite3{queries: queries, db: db}, err
}

func (s *SQLite3) Close() error {
	return s.db.Close()
}

func (s *SQLite3) CreateMessage(
	ctx context.Context,
	params sqlc.CreateMessageParams,
) (sqlc.Message, error) {
	return s.queries.CreateMessage(ctx, params)
}

func (s *SQLite3) CreateChat(ctx context.Context, name sql.NullString) (sqlc.Chat, error) {
	return s.queries.CreateChat(ctx, name)
}

func (s *SQLite3) CreateChatAndMessage(
	ctx context.Context,
	params sqlc.CreateMessageParams,
) (sqlc.Chat, sqlc.Message, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return sqlc.Chat{}, sqlc.Message{}, err
	}
	defer tx.Rollback()
	qtx := s.queries.WithTx(tx)
	chat, err := qtx.CreateChat(ctx, sql.NullString{String: params.Content})
	if err != nil {
		return sqlc.Chat{}, sqlc.Message{}, err
	}
	params.ChatID = chat.ID
	msg, err := qtx.CreateMessage(ctx, params)
	if err != nil {
		return sqlc.Chat{}, sqlc.Message{}, err
	}
	err = tx.Commit()
	return chat, msg, err
}
