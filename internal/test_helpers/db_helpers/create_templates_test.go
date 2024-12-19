package db_helpers_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestCreateTemplate(t *testing.T) {
	db := prepare(t)
	ctx := context.Background()
	expected := models.Template{
		Name:    "hello world",
		Content: `--{{block "hello world"}}hello world: --{{.}}--{{end}}`,
	}
	err := db_helpers.CreateTemplate(db, ctx, &expected)
	if err != nil {
		t.Fatalf("Error upon templates creation: %s\n", err)
	}
	row := db.QueryRow("SELECT name, content FROM templates WHERE ID = 1")
	if err := row.Err(); err != nil {
		t.Fatalf("Error upon quering template: %s\n", err)
	}
	var actual models.Template
	if err := row.Scan(&actual.Name, &actual.Content); err != nil {
		t.Fatalf("Failed to scan rows, errored: %s\n", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf(
			"Expected and actual results are different\nexpected: %+v\nactual:   %+v\n",
			expected,
			actual,
		)
	}
}

func TestGenerateTemplatesSliceN(t *testing.T) {
	test_helpers.Skip(t)
	messages := db_helpers.GenerateTemplateSliceN(3)
	actual := test_helpers.UnpointerSlice(messages)
	expected := []models.Template{
		{ID: 1, Content: `--{{template "1"}}1--{{end}}`, Name: "1"},
		{ID: 2, Content: `--{{template "2"}}2--{{end}}`, Name: "2"},
		{ID: 3, Content: `--{{template "3"}}3--{{end}}`, Name: "3"},
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"Bad result in generating messages slice\nexpected: %+v\nactual:   %+v\n",
			expected,
			actual,
		)
	}
}
