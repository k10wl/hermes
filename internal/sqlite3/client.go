package sqlite3

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLite3 struct {
	DB *sql.DB
}

func NewSQLite3(dsn string) (*SQLite3, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err := runMigrations(db); err != nil {
		return nil, err
	}
	if err := writePragma(db); err != nil {
		return nil, err
	}
	return &SQLite3{DB: db}, err
}

func (s *SQLite3) Close() error {
	return s.DB.Close()
}

func writePragma(db *sql.DB) error {
	_, err := db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return fmt.Errorf("Unable to set params: %s", err)
	}
	return nil
}
