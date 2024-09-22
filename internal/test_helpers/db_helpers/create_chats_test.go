package db_helpers_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestCreateChats(t *testing.T) {
	db := prepare(t)
	defer db.Close()
	ctx := context.Background()

	subject := []*models.Chat{}
	for i := 0; i < 10; i++ {
		subject = append(subject, &models.Chat{Name: strconv.Itoa(i)})
	}

	err := db_helpers.CreateChats(db, ctx, subject)
	if err != nil {
		t.Errorf("error upon chat creation - %s\n", err)
		return
	}

	rows, err := db.Query("SELECT name FROM chats")
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
