package db_helpers_test

import (
	"context"
	"testing"

	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestChatsSeeder(t *testing.T) {
	db := prepare(t)
	defer db.Close()
	ctx := context.Background()

	seeder := db_helpers.NewSeeder(db, ctx)
	seedsAmount := 420
	seeder.SeedChatsN(int64(seedsAmount))

	rows, err := db.Query(`SELECT * FROM chats`)
	if err != nil {
		t.Errorf("error upon chats query - %s\n", err)
	}
	for rows.Next() {
		seedsAmount--
	}
	if seedsAmount != 0 {
		t.Errorf(
			"amount of seeds did not match amount of rows, mismatch: %d\n",
			seedsAmount,
		)
	}
}

func TestMessagesSeeder(t *testing.T) {
	db := prepare(t)
	defer db.Close()
	ctx := context.Background()

	seeder := db_helpers.NewSeeder(db, ctx)
	seedsAmount := 420
	seeder.SeedMessagesN(int64(seedsAmount), 1)

	rows, err := db.Query(`SELECT * FROM messages`)
	if err != nil {
		t.Errorf("error upon messages query - %s\n", err)
	}
	for rows.Next() {
		seedsAmount--
	}
	if seedsAmount != 0 {
		t.Errorf(
			"amount of seeds did not match amount of rows, mismatch: %d\n",
			seedsAmount,
		)
	}
}
