package sqlite3

import (
	"database/sql"

	"github.com/k10wl/hermes/internal/settings"
	_ "modernc.org/sqlite"
)

type SQLite3 struct {
	db *sql.DB
}

func NewSQLite3(config *settings.Config) (*SQLite3, error) {
	dbName := config.DatabaseDSN
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	err = runMigrations(db)
	if err != nil {
		return nil, err
	}
	return &SQLite3{db: db}, err
}

func (s *SQLite3) Close() error {
	return s.db.Close()
}
