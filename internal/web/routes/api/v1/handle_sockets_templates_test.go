package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
)

func TestGetTemplates(t *testing.T) {
	client, db, teardown := setupWebSocketTest(t)
	defer teardown()

	seeder := db_helpers.NewSeeder(db, context.Background())
	seeded, err := seeder.SeedTemplatesN(100)
	if err != nil {
		t.Fatalf("error upon seeding templates - %s\n", err)
	}
	slices.Reverse(seeded)

	if err := client.WriteMessage(
		websocket.TextMessage,
		[]byte(fmt.Sprintf(
			`{"id": %q, "type": "request-read-templates", "payload": {"start_before_id": -1, "limit": -1}}`,
			sharedID,
		)),
	); err != nil {
		t.Fatalf("Failed to send request read templates message from client, error: %v\n", err)
	}

	_, data, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read request templates message: %v\n", err)
	}

	resData := new(messages.ServerReadTemplates)
	if err := json.Unmarshal(data, resData); err != nil {
		t.Fatalf("failed to unmarshal response - %s\n", err)
	}

	actual := test_helpers.UnpointerSlice(test_helpers.ResetSliceTime(resData.Payload.Templates))
	expected := test_helpers.UnpointerSlice(test_helpers.ResetSliceTime(seeded))
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf(
			"Failed to read templates, results differ from expected.\nexpected: %+v\nactual:   %+v\n",
			expected,
			actual,
		)
	}
}

func TestGetTemplatesWithName(t *testing.T) {
	client, db, teardown := setupWebSocketTest(t)
	defer teardown()

	seeder := db_helpers.NewSeeder(db, context.Background())
	_, err := seeder.SeedTemplatesN(100)
	if err != nil {
		t.Fatalf("error upon seeding templates - %s\n", err)
	}

	if err := client.WriteMessage(
		websocket.TextMessage,
		[]byte(fmt.Sprintf(
			`{"id": %q, "type": "request-read-templates", "payload": {"start_before_id": -1, "limit": -1, "name": "22"}}`,
			sharedID,
		)),
	); err != nil {
		t.Fatalf("Failed to send request read templates message from client, error: %v\n", err)
	}

	_, data, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read request templates message: %v\n", err)
	}

	resData := new(messages.ServerReadTemplates)
	if err := json.Unmarshal(data, resData); err != nil {
		t.Fatalf("failed to unmarshal response - %s\n", err)
	}

	actual := test_helpers.UnpointerSlice(
		test_helpers.ResetSliceTime(resData.Payload.Templates),
	)
	expected := []models.Template{
		{
			ID:      22,
			Name:    "22",
			Content: `--{{template "22"}}22--{{end}}`,
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf(
			"Failed to read templates, results differ from expected.\nexpected: %+v\nactual:   %+v\n",
			expected,
			actual,
		)
	}
}

func TestGetTemplateByID(t *testing.T) {
	client, db, teardown := setupWebSocketTest(t)
	defer teardown()

	seeder := db_helpers.NewSeeder(db, context.Background())
	templates, err := seeder.SeedTemplatesN(10)
	if err != nil {
		t.Fatalf("error upon seeding templates - %s\n", err)
	}

	id := 2
	if err := client.WriteMessage(
		websocket.TextMessage,
		[]byte(fmt.Sprintf(
			`{"id": %q, "type": "request-read-template", "payload": {"id": %d}}`,
			sharedID,
			id,
		)),
	); err != nil {
		t.Fatalf("Failed to send request read template message from client, error: %v\n", err)
	}

	_, data, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read request template message: %v\n", err)
	}

	resData := new(messages.ServerReadTemplate)
	if err := json.Unmarshal(data, resData); err != nil {
		t.Fatalf("failed to unmarshal response - %s\ndata: %s\n", err, data)
	}

	resData.Payload.Template.TimestampsToNilForTest__()
	actual := *resData.Payload.Template
	expected := *templates[id-1]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf(
			"Failed to read template, results differ from expected.\nexpected: %+v\nactual:   %+v\n",
			expected,
			actual,
		)
	}
}

func TestGetTemplateByIDErrorsUponNonExisting(t *testing.T) {
	client, _, teardown := setupWebSocketTest(t)
	defer teardown()

	if err := client.WriteMessage(
		websocket.TextMessage,
		[]byte(fmt.Sprintf(
			`{"id": %q, "type": "request-read-template", "payload": {"id": 99999}}`,
			sharedID,
		)),
	); err != nil {
		t.Fatalf("Failed to send request read template message from client, error: %v\n", err)
	}

	_, data, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read request template message: %v\n", err)
	}

	resData := new(messages.ServerError)
	if err := json.Unmarshal(data, resData); err != nil {
		t.Fatalf(
			"Expected error, but got something else\nUnmarshal error: %s\n data: %s\n",
			err.Error(),
			data,
		)
	}
}
