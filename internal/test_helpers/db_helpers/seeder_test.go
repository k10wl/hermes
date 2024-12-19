package db_helpers_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
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
	if err := seeder.SeedChatsN(1); err != nil {
		t.Fatalf("error upon seeding messages - %s\n", err)
	}
	if err := seeder.SeedMessagesN(int64(seedsAmount), 1); err != nil {
		t.Fatalf("error upon seeding messages - %s\n", err)
	}

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

func TestTemplatesSeeder(t *testing.T) {
	db := prepare(t)
	defer db.Close()
	ctx := context.Background()

	seeder := db_helpers.NewSeeder(db, ctx)
	seedsAmount := 3

	seeded, err := seeder.SeedTemplatesN(int64(seedsAmount))
	if err != nil {
		t.Fatalf("error upon seeding templates - %s\n", err)
	}

	rows, err := db.Query(`SELECT * FROM templates`)
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

	expected := []models.Template{
		{ID: 1, Content: `--{{template "1"}}1--{{end}}`, Name: "1"},
		{ID: 2, Content: `--{{template "2"}}2--{{end}}`, Name: "2"},
		{ID: 3, Content: `--{{template "3"}}3--{{end}}`, Name: "3"},
	}
	actual := test_helpers.UnpointerSlice(seeded)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"Bad result in generating messages slice\nexpected: %+v\nactual:   %+v\n",
			expected,
			actual,
		)
	}
}
