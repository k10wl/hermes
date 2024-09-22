package db_helpers_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/sqlite3"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestCreateChats(t *testing.T) {
	test_helpers.Skip(t)
	sqlite3, err := sqlite3.NewSQLite3(":memory:")
	if err != nil {
		t.Errorf("error during setup - %s\n", err)
		return
	}

	subject := []*models.Chat{}
	for i := 0; i < 10; i++ {
		subject = append(subject, &models.Chat{Name: strconv.Itoa(i)})
	}

	err = db_helpers.CreateChats(sqlite3.DB, context.Background(), subject)
	if err != nil {
		t.Errorf("error upon chat creation - %s\n", err)
		return
	}

	rows, err := sqlite3.DB.Query("SELECT name FROM chats")
	if err != nil {
		t.Errorf("error upon db query - %s\n", err)
		return
	}

	for i, subject := range subject {
		ok := rows.Next()
		if !ok {
			t.Errorf("unexpected end of rows on i = %d\n", i)
			break
		}
		record := models.Chat{}
		err := rows.Scan(&record.Name)
		if err != nil {
			t.Errorf("scanning error - %s\n", err)
			break
		}
		if subject.Name != record.Name {
			t.Errorf(
				"bad name\nexpected: %s\nactual:   %s\n",
				subject.Name,
				record.Name,
			)
			break
		}
	}
}
