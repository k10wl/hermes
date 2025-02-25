package db_helpers_test

import (
	"context"
	"testing"

	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestGetChatByID(t *testing.T) {
	db := prepare(t)
	defer db.Close()
	ctx := context.Background()

	rows, err := db.ExecContext(ctx, `
INSERT INTO
    chats (id, name)
VALUES 
    (1, 'first'),
    (2, 'second')
`)
	if err != nil {
		t.Errorf("failed to insert into chats - %s\n", err)
		return
	}
	n, err := rows.RowsAffected()
	if err != nil {
		t.Errorf("failed to count inserted rows - %s\n", err)
		return
	}
	if n == 0 {
		t.Errorf("no rows affected upon insertion\n")
		return
	}

	chat, err := db_helpers.GetChatByID(db, ctx, 1)
	if err != nil {
		t.Errorf("error upon get chat by id - %s\n", err)
		return
	}
	if chat.Name != "first" || chat.ID != 1 {
		t.Errorf("retracted wrong chat, expected first - %+v\n", chat)
		return
	}

	chat, err = db_helpers.GetChatByID(db, ctx, 2)
	if err != nil {
		t.Errorf("error upon get chat by id - %s\n", err)
		return
	}
	if chat.Name != "second" || chat.ID != 2 {
		t.Errorf("retracted wrong chat, expected second - %+v\n", chat)
		return
	}

	_, err = db_helpers.GetChatByID(db, ctx, 3)
	if err == nil {
		t.Errorf("expected to error upon non existing chat but received nil\n")
		return
	}
}
