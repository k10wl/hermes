package sqlite3

import (
	"database/sql"

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
	err = runMigrations(db)
	if err != nil {
		return nil, err
	}
	return &SQLite3{DB: db}, err
}

func (s *SQLite3) Close() error {
	return s.DB.Close()
}
