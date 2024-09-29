package db_helpers_test

import (
	"database/sql"
	"testing"

	"github.com/k10wl/hermes/internal/sqlite3"
	"github.com/k10wl/hermes/internal/test_helpers"
)

func prepare(t *testing.T) *sql.DB {
	test_helpers.Skip(t)
	sqlite3, err := sqlite3.NewSQLite3(":memory:")
	if err != nil {
		t.Fatalf("failed to setup database - %s\n", err)
	}
	return sqlite3.DB
}
