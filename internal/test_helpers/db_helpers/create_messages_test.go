package db_helpers_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestCreateMessages(t *testing.T) {
	db := prepare(t)
	defer db.Close()
	ctx := context.Background()

	subject := []*models.Message{}
	for i := 0; i < 10; i++ {
		subject = append(
			subject,
			&models.Message{
				ChatID:  1,
				Content: "generated",
				Role:    "assistant",
				ID:      int64(i),
			},
		)
	}

	err := db_helpers.CreateMessages(db, ctx, subject)
	if err != nil {
		t.Errorf("error upon messages creation - %s\n", err)
		return
	}

	rows, err := db.Query("SELECT content FROM messages WHERE chat_id = 1")
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
		record := models.Message{}
		err := rows.Scan(&record.Content)
		if err != nil {
			t.Errorf("scanning error - %s\n", err)
			break
		}
		if subject.Content != record.Content {
			t.Errorf(
				"bad name\nexpected: %s\nactual:   %s\n",
				subject.Content,
				record.Content,
			)
			break
		}
	}
}

func TestGenerateMessagesSliceN(t *testing.T) {
	test_helpers.Skip(t)
	messages := db_helpers.GenerateMessagesSliceN(3, 1)
	actual := test_helpers.UnpointerSlice(messages)
	expected := []models.Message{
		{ID: 1, Content: "generated", Role: "assistant", ChatID: 1},
		{ID: 2, Content: "generated", Role: "assistant", ChatID: 1},
		{ID: 3, Content: "generated", Role: "assistant", ChatID: 1},
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"Bad result in generating messages slice\nexpected: %+v\nactual:   %+v\n",
			expected,
			actual,
		)
	}
}
