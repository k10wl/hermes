package db_helpers_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestFindTeplateByName(t *testing.T) {
	db := prepare(t)

	if _, err := db_helpers.FindTemplateByName(
		db,
		context.Background(),
		"first",
	); err == nil {
		t.Fatal("Failed to report missing template\n")
	}

	if _, err := db.Exec(`
INSERT INTO templates (name, content)
VALUES ('first', '--{{define "first"}}[--{{.}}]--{{end}}');
`); err != nil {
		t.Fatalf("Failed to insert test template - %v\n", err)
	}

	template, err := db_helpers.FindTemplateByName(db, context.Background(), "first")
	if err != nil {
		t.Fatalf("Unexpected error in FindTemplateByName- %v\n", err)
	}
	template.TimestampsToNilForTest__()
	expected := models.Template{
		Name:    "first",
		Content: "--{{define \"first\"}}[--{{.}}]--{{end}}",
		ID:      1,
	}
	if !reflect.DeepEqual(*template, expected) {
		t.Fatalf(
			"Failed to find correct template\nexpected: %+v\nactual:   %+v\n",
			expected,
			*template,
		)
	}
}
