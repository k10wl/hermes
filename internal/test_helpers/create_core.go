package test_helpers

import (
	"database/sql"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
)

func CreateCore() (*core.Core, *sql.DB) {
	db, err := sqlite3.NewSQLite3(":memory:")
	if err != nil {
		panic(err)
	}
	return core.NewCore(db, &settings.Config{}), db.DB
}
