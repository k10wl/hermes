package test_helpers

import (
	"context"
	"database/sql"
	"strings"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
)

func CreateCore() (*core.Core, *sql.DB) {
	db, err := sqlite3.NewSQLite3(":memory:")
	if err != nil {
		panic(err)
	}
	c := core.NewCore(db, &settings.Config{})
	c.GetConfig().Stdoout = &strings.Builder{}
	c.GetConfig().Stderr = &strings.Builder{}
	c.GetConfig().ShutdownContext = context.Background()
	return c, db.DB
}
