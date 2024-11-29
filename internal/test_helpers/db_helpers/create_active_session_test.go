package db_helpers_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestActiveSession(t *testing.T) {
	db := prepare(t)
	expected := models.ActiveSession{
		ID:          1,
		Address:     "42069",
		DatabaseDNS: "deez",
	}
	db_helpers.CreateActiveSession(db, context.Background(), &expected)
	row := db.QueryRow(
		`SELECT id, address, database_dns FROM active_sessions WHERE id = 1`,
	)
	if err := row.Err(); err != nil {
		t.Fatalf("error upon active sessions query: %s\n", err)
	}
	var actual models.ActiveSession
	if err := row.Scan(
		&actual.ID,
		&actual.Address,
		&actual.DatabaseDNS,
	); err != nil {
		t.Fatalf("error upon scanning active session result: %s\n", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf(
			"failed to get expected active session\nexpected: %+v\nactual:   %+v\n",
			expected,
			actual,
		)
	}
}
