package sqlite3

import (
	"path"

	"github.com/k10wl/hermes/internal/runtime"
)

type SQLite3 struct{}

func NewSQLite3(config *runtime.Config) (*SQLite3, error) {
	err := runMigrations(path.Join(config.ConfigDir, "main.db"))
	return &SQLite3{}, err
}
