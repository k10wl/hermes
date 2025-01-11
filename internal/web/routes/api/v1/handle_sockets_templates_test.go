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
			Content: `--{{define "22"}}22--{{end}}`,
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

func TestEditTemplateContent(t *testing.T) {
	client, db, teardown := setupWebSocketTest(t)
	defer teardown()
	seeder := db_helpers.NewSeeder(db, context.Background())

	if _, err := seeder.SeedTemplatesN(1); err != nil {
		t.Fatalf("seeding templates failed - %q\n", err)
	}

	editContent := `--{{define "1"}}edited--{{end}}`

	if err := client.WriteMessage(
		websocket.TextMessage,
		[]byte(fmt.Sprintf(
			`
{
  "id": %q,
  "type": "request-edit-template",
  "payload": {
    "name": "1",
    "content": %q
  }
}
`,
			sharedID,
			editContent,
		)),
	); err != nil {
		t.Fatalf("Failed to send request read template message from client, error: %v\n", err)
	}

	_, bytes, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message - %q\n", err)
	}

	templateChanged := new(messages.ServerTemplateChanged)
	json.Unmarshal(bytes, templateChanged)

	if templateChanged.Payload.Template.Content != editContent {
		t.Fatalf(
			"Failed to retrieve edited content\nexpected: %q\nactual:   %q\n",
			editContent,
			templateChanged.Payload.Template.Content,
		)
	}
}

func TestEditTemplateRename(t *testing.T) {
	client, db, teardown := setupWebSocketTest(t)
	defer teardown()
	ctx := context.Background()
	seeder := db_helpers.NewSeeder(db, ctx)

	if _, err := seeder.SeedTemplatesN(1); err != nil {
		t.Fatalf("seeding templates failed - %q\n", err)
	}

	editContent := `--{{define "edited"}}edited--{{end}}`

	if err := client.WriteMessage(
		websocket.TextMessage,
		[]byte(fmt.Sprintf(
			`
{
  "id": %q,
  "type": "request-edit-template",
  "payload": {
    "name": "1",
    "content": %q
  }
}
`,
			sharedID,
			editContent,
		)),
	); err != nil {
		t.Fatalf("Failed to send request read template message from client, error: %v\n", err)
	}

	_, bytes, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message - %q\n", err)
	}

	templateChanged := new(messages.ServerTemplateChanged)
	json.Unmarshal(bytes, templateChanged)
	templateChanged.Payload.Template.TimestampsToNilForTest__()

	if templateChanged.Payload.Template.Content != editContent {
		t.Fatalf(
			"Failed to retrieve edited content\nexpected: %q\nactual:   %q\n",
			editContent,
			templateChanged.Payload.Template.Content,
		)
	}

	if _, err := db_helpers.FindTemplateByName(db, ctx, "1"); err == nil {
		t.Fatalf(
			"Template with previous name still exists\n",
		)
	}

	dbTemplate, err := db_helpers.FindTemplateByName(db, ctx, "edited")
	if err != nil {
		t.Fatalf("Failed to find new template name current template\n")
	}
	dbTemplate.TimestampsToNilForTest__()

	if !reflect.DeepEqual(*dbTemplate, templateChanged.Payload.Template) {
		t.Fatalf(
			"Failed to read template, results differ from expected.\nexpected: %+v\nactual:   %+v\n",
			*dbTemplate,
			templateChanged.Payload.Template,
		)
	}
}

func TestEditTemplateRenameWithClone(t *testing.T) {
	client, db, teardown := setupWebSocketTest(t)
	defer teardown()
	ctx := context.Background()
	seeder := db_helpers.NewSeeder(db, ctx)

	if _, err := seeder.SeedTemplatesN(1); err != nil {
		t.Fatalf("seeding templates failed - %q\n", err)
	}

	editContent := `--{{define "edited"}}edited--{{end}}`

	if err := client.WriteMessage(
		websocket.TextMessage,
		[]byte(fmt.Sprintf(
			`
{
  "id": %q,
  "type": "request-edit-template",
  "payload": {
    "name": "1",
    "content": %q,
    "clone": true
  }
}
`,
			sharedID,
			editContent,
		)),
	); err != nil {
		t.Fatalf("Failed to send request read template message from client, error: %v\n", err)
	}

	_, bytes, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message - %q\n", err)
	}

	templateCreated := new(messages.ServerTemplateCreated)
	json.Unmarshal(bytes, templateCreated)
	templateCreated.Payload.Template.TimestampsToNilForTest__()

	if templateCreated.Payload.Template.Content != editContent {
		t.Fatalf(
			"Failed to retrieve edited content\nexpected: %q\nactual:   %q\n",
			editContent,
			templateCreated.Payload.Template.Content,
		)
	}
	tmp, err := db_helpers.FindTemplateByName(db, ctx, "1")
	if err != nil {
		t.Fatalf(
			"Failed to get initial template name, error: %s\n", err,
		)
	}
	initialExpected := `--{{define "1"}}1--{{end}}`
	if tmp.Content != initialExpected {
		t.Fatalf(
			"Initial template content differs\nexpected: %q\nactual:   %q\n",
			tmp.Content,
			initialExpected,
		)
	}

	dbTemplate, err := db_helpers.FindTemplateByName(db, ctx, "edited")
	if err != nil {
		t.Fatalf("Failed to find new template name current template\n")
	}
	dbTemplate.TimestampsToNilForTest__()

	if !reflect.DeepEqual(*dbTemplate, templateCreated.Payload.Template) {
		t.Fatalf(
			"Failed to read template, results differ from expected.\nexpected: %+v\nactual:   %+v\n",
			*dbTemplate,
			templateCreated.Payload.Template,
		)
	}
}
