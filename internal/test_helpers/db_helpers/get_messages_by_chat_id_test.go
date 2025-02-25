package db_helpers_test

import (
	"context"
	"testing"

	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestGetMessagesByChatID(t *testing.T) {
	db := prepare(t)
	q := `
INSERT INTO chats (name) VALUES ('generic');
INSERT INTO messages (chat_id, role_id, content) VALUES 
    (1, 1, 'user question'),
    (1, 2, 'assistant response');
    `
	_, err := db.Exec(q)
	if err != nil {
		t.Fatalf("Failed to prepare db for test, error: %s\n", err)
	}
	messages, err := db_helpers.GetMessagesByChatID(db, context.Background(), 1)
	if err != nil {
		t.Fatalf("Error upon getting messages by chat id, error: %s\n", err)
	}
	if len(messages) != 2 {
		t.Fatalf(
			"Incorrect query length, db has 2 records, but function returned %d\n",
			len(messages),
		)
	}
}
