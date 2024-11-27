package core_test

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestCreateChatAndCompletionCommand(t *testing.T) {
	type testCase struct {
		name           string
		init           func()
		shouldError    bool
		expectedResult models.Message
	}

	coreInstance, _ := test_helpers.CreateCore()
	var currentCommand *core.CreateChatAndCompletionCommand

	dbTemplates := map[string]string{
		"welcome": `--{{define "welcome"}}hello world!--{{end}}`,
		"wrapper": `--{{define "wrapper"}}wrapper - --{{.}} - wrapper--{{end}}`,
		"nested1": `--{{define "nested1"}}nested1: --{{template "nested2" .}}--{{end}}`,
		"nested2": `--{{define "nested2"}}nested2: --{{.}}--{{end}}`,
		"loop1":   `--{{define "loop1"}}--{{template "loop2"}}--{{end}}`,
		"loop2":   `--{{define "loop2"}}--{{template "loop1"}}--{{end}}`,
	}

	for _, template := range dbTemplates {
		if err := core.NewUpsertTemplateCommand(
			coreInstance,
			template,
		).Execute(context.Background()); err != nil {
			panic(err)
		}
	}

	table := []testCase{
		{
			name: "Should create simple response",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					"hello world",
					"",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      2,
				ChatID:  1,
				Role:    core.AssistantRole,
				Content: "> mocked: hello world",
			},
		},
		{
			name: "Should fill template data when template is passed as argument",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					``,
					"welcome",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      4,
				ChatID:  2,
				Role:    core.AssistantRole,
				Content: "> mocked: hello world!",
			},
		},
		{
			name: "Should fill inner template data when template name is not provided",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					`--{{template "welcome"}}`,
					"",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      6,
				ChatID:  3,
				Role:    core.AssistantRole,
				Content: "> mocked: hello world!",
			},
		},
		{
			name: "Should fill inner template data and provided template name data",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					`--{{template "welcome"}}`,
					"wrapper",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      8,
				ChatID:  4,
				Role:    core.AssistantRole,
				Content: "> mocked: wrapper - hello world! - wrapper",
			},
		},
		{
			name: "Should error if given template name does not exist",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					`--{{template "welcome"}}`,
					"does not exist",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError:    true,
			expectedResult: models.Message{},
		},
		{
			name: "Should process string input and fill template data",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					"hello world!",
					"wrapper",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      10,
				ChatID:  5,
				Role:    core.AssistantRole,
				Content: "> mocked: wrapper - hello world! - wrapper",
			},
		},
		{
			name: "Should error on circular templates with provided template",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					"will blow up",
					"loop1",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError:    true,
			expectedResult: models.Message{},
		},
		{
			name: "Should error on circular templates with inputted template",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					`will blow up --{{template "loop1" . }} `,
					"",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError:    true,
			expectedResult: models.Message{},
		},
		{
			name: "Should error on circular templates with inputted template and template name as an argument",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					`will blow up --{{template "loop1" . }} `,
					"loop2",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError:    true,
			expectedResult: models.Message{},
		},
		{
			name: "Should fill inner templates and remain current message",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					`should fill welcome (--{{template "welcome"}})(--{{template "welcome"}})`,
					"",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ChatID:  6,
				ID:      12,
				Role:    core.AssistantRole,
				Content: "> mocked: should fill welcome (hello world!)(hello world!)",
			},
		},
		{
			name: "Should fill inner templates, template name, and remain current message",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance,
					core.AssistantRole,
					`should fill welcome (--{{template "welcome"}})(--{{template "welcome"}})`,
					"wrapper",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ChatID:  7,
				ID:      14,
				Role:    core.AssistantRole,
				Content: "> mocked: wrapper - should fill welcome (hello world!)(hello world!) - wrapper",
			},
		},
	}

	for _, test := range table {
		test.init()
		err := currentCommand.Execute(context.TODO())
		test.expectedResult.TimestampsToNilForTest__()
		if currentCommand.Result != nil {
			currentCommand.Result.TimestampsToNilForTest__()
		}
		if test.shouldError {
			if err == nil {
				t.Errorf(
					"%q expected to error, but did not.\nres: %+v\n\n",
					test.name,
					*currentCommand.Result,
				)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q unexpected error: %v\n\n", test.name, err)
			continue
		}
		if !reflect.DeepEqual(test.expectedResult, *currentCommand.Result) {
			t.Errorf(
				"%q - bad result\nexpected: %+v\nactual:   %+v\n\n",
				test.name,
				test.expectedResult,
				*currentCommand.Result,
			)
		}
	}
}

func TestCreateCompletionCommand(t *testing.T) {
	type testCase struct {
		name           string
		init           func()
		shouldError    bool
		expectedResult models.Message
	}

	coreInstance, db := test_helpers.CreateCore()
	err := db_helpers.NewSeeder(db, context.Background()).SeedChatsN(1)
	if err != nil {
		t.Fatalf("Failed to seed chats: %s\n", err)
	}
	var currentCommand *core.CreateCompletionCommand

	table := []testCase{
		{
			name: "Should create simple response",
			init: func() {
				currentCommand = core.NewCreateCompletionCommand(
					coreInstance,
					1,
					core.AssistantRole,
					"hello world",
					"",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      2,
				ChatID:  1,
				Role:    core.AssistantRole,
				Content: "> mocked: hello world",
			},
		},
	}

	for _, test := range table {
		test.init()
		err := currentCommand.Execute(context.TODO())
		test.expectedResult.TimestampsToNilForTest__()
		currentCommand.Result.TimestampsToNilForTest__()
		if test.shouldError && err == nil {
			t.Errorf("%s expected to error, but did not\n", test.name)
			continue
		}
		if !test.shouldError && err != nil {
			t.Errorf("%s unexpected error: %v\n", test.name, err)
			continue
		}
		if !reflect.DeepEqual(test.expectedResult, *currentCommand.Result) {
			t.Errorf(
				"%s - bad result\nexpected: %+v\nactual:   %+v",
				test.name,
				test.expectedResult,
				*currentCommand.Result,
			)
		}
	}
}

func TestCreateTemplateCommand(t *testing.T) {
	type testCase struct {
		name         string
		template     string
		templateName []string
		init         func()
		shouldError  bool
	}

	coreInstance, _ := test_helpers.CreateCore()
	db := coreInstance.GetDB()
	var command *core.UpsertTemplateCommand

	table := []testCase{
		{
			name:         "create welcome template",
			template:     `--{{define "welcome"}}hello world!--{{end}}`,
			templateName: []string{"welcome"},
			init: func() {
				command = core.NewUpsertTemplateCommand(
					coreInstance,
					`--{{define "welcome"}}hello world!--{{end}}`,
				)
			},
			shouldError: false,
		},
		{
			name:         "should override written command",
			template:     `--{{define "welcome"}}welcome world!--{{end}}`,
			templateName: []string{"welcome"},
			init: func() {
				command = core.NewUpsertTemplateCommand(
					coreInstance,
					`--{{define "welcome"}}welcome world!--{{end}}`,
				)
			},
			shouldError: false,
		},
	}

	for _, test := range table {
		test.init()
		err := command.Execute(context.TODO())
		if test.shouldError && err == nil {
			t.Errorf("%q expected to error, but did not\n", test.name)
			continue
		}
		if !test.shouldError && err != nil {
			t.Errorf("%q unexpected error: %v\n", test.name, err)
			continue
		}
		tmpl, err := db.GetTemplatesByNames(context.TODO(), test.templateName)
		if err != nil {
			t.Errorf("%q failed to get created template: %v\n", test.name, err)
			continue
		}
		if test.template != tmpl[0].Content {
			t.Errorf(
				"%q - bad result\nexpected: %+v\nactual:   %+v\n",
				test.name,
				test.template,
				tmpl,
			)
		}
	}
}

func TestDeleteTemplateByName(t *testing.T) {
	type testCase struct {
		name         string
		init         func()
		shouldError  bool
		templateName string
	}

	coreInstance, _ := test_helpers.CreateCore()
	cmd := core.NewDeleteTemplateByName
	var command core.DeleteTemplateByName

	templates := []string{
		`--{{define "welcome1"}}welcome{{end}}`,
		`--{{define "welcome2"}}welcome{{end}}`,
	}

	for _, template := range templates {
		upsertCmd := core.NewUpsertTemplateCommand(coreInstance, template)
		if err := upsertCmd.Execute(context.Background()); err != nil {
			panic("bad test setup")
		}
	}

	table := []testCase{
		{
			name: "should delete template",
			init: func() {
				command = *cmd(coreInstance, "welcome1")
			},
			templateName: "welcome1",
		},
		{
			name: "should return err if no templates were deleted",
			init: func() {
				command = *cmd(coreInstance, "does not exist")
			},
			templateName: "does not exist",
			shouldError:  true,
		},
		{
			name: "should delete second template",
			init: func() {
				command = *cmd(coreInstance, "welcome2")
			},
			templateName: "welcome2",
		},
	}

	for _, test := range table {
		test.init()
		err := command.Execute(context.Background())
		if test.shouldError && err == nil {
			t.Errorf("%q expected to error, but did not\n\n", test.name)
			continue
		}
		if !test.shouldError && err != nil {
			t.Errorf("%q unexpected error: %v\n\n", test.name, err)
			continue
		}
		getTemplates := core.NewGetTemplatesByNamesQuery(coreInstance, []string{test.templateName})
		if err := getTemplates.Execute(context.Background()); err != nil {
			t.Errorf("%q query templates error: %v\n\n", test.name, err)
		}
		if len(getTemplates.Result) != 0 {
			t.Errorf("%q failed to delete template\n\n", test.name)
			continue
		}
	}
}

func TestEditTemplateByName(t *testing.T) {
	type testCase struct {
		name                string
		init                func() string
		templateName        string
		shouldError         bool
		deletedTemplateName string
	}

	coreInstance, _ := test_helpers.CreateCore()
	var command core.EditTemplateByName

	templates := []string{
		`--{{define "welcome"}}welcome--{{end}}`,
	}
	for _, template := range templates {
		upsertCmd := core.NewUpsertTemplateCommand(coreInstance, template)
		if err := upsertCmd.Execute(context.Background()); err != nil {
			panic("bad test setup")
		}
	}

	table := []testCase{
		{
			name: "should edit existing template",
			init: func() string {
				content := `--{{define "welcome"}}hi--{{end}}`
				command = *core.NewEditTemplateByName(coreInstance, "welcome", content, false)
				return content
			},
			templateName: "welcome",
		},
		{
			name: "should allow block definition",
			init: func() string {
				content := `--{{block "welcome" .}}hi--{{end}}`
				command = *core.NewEditTemplateByName(coreInstance, "welcome", content, false)
				return content
			},
			templateName: "welcome",
		},
		{
			name: "should error if template does not exist",
			init: func() string {
				content := `--{{define "welcome"}}hi--{{end}}`
				command = *core.NewEditTemplateByName(coreInstance, "does not exist", content, false)
				return content
			},
			templateName: "does not exist",
			shouldError:  true,
		},
		{
			name: "should error if template has corrupted content",
			init: func() string {
				content := `--{{define "welco`
				command = *core.NewEditTemplateByName(coreInstance, "welcome", content, false)
				return content
			},
			templateName: "welcome",
			shouldError:  true,
		},
		{
			name: "should rename template if new name does not match existing name and clone is false",
			init: func() string {
				content := `--{{define "missmatch"}}stuff--{{end}}`
				command = *core.NewEditTemplateByName(coreInstance, "welcome", content, false)
				return content
			},
			templateName:        "missmatch",
			deletedTemplateName: "welcome",
		},
		{
			name: "should create new template if cloning is true and new name is unique",
			init: func() string {
				content := `--{{define "missmatch2"}}stuff--{{end}}`
				command = *core.NewEditTemplateByName(coreInstance, "welcome", content, true)
				return content
			},
			templateName: "missmatch2",
		},
		{
			name: "should fail if cloning is true but new name is not unique",
			init: func() string {
				content := `--{{define "missmatch2"}}stuff--{{end}}`
				command = *core.NewEditTemplateByName(coreInstance, "missmatch", content, true)
				return content
			},
			templateName: "missmatch2",
			shouldError:  true,
		},
	}

	for _, test := range table {
		expected := test.init()
		err := command.Execute(context.Background())
		if test.shouldError {
			if err == nil {
				t.Errorf("%q expected to error, but did not\n\n", test.name)
			}
			continue
		}
		if !test.shouldError && err != nil {
			t.Errorf("%q unexpected error: %v\n\n", test.name, err)
			continue
		}
		getTemplates := core.NewGetTemplatesByNamesQuery(coreInstance, []string{test.templateName})
		if err := getTemplates.Execute(context.Background()); err != nil {
			t.Errorf("%q query templates error: %v\n\n", test.name, err)
		}
		if len(getTemplates.Result) != 1 {
			t.Errorf("%q failed to edit template\n\n", test.name)
			continue
		}
		actual := getTemplates.Result[0].Content
		if expected != actual {
			t.Errorf(
				"%q bad result\nexpected: %v\nactual:   %v\n\n",
				test.name,
				expected,
				actual,
			)
			continue
		}
		if test.deletedTemplateName != "" {
			query := core.NewGetTemplatesByNamesQuery(
				coreInstance,
				[]string{test.deletedTemplateName},
			)
			if len(query.Result) != 0 {
				t.Errorf(
					"%q did not remove original template - %q",
					test.name,
					test.deletedTemplateName,
				)
				continue
			}
		}
	}
}

func TestCreateChat(t *testing.T) {
	c, db := test_helpers.CreateCore()
	cmd := core.NewCreateChatWithMessageCommand(c, &models.Message{
		Role:    "user",
		Content: "hello-world",
	})
	err := cmd.Execute(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	dbChat, err := db_helpers.GetChatByID(db, context.Background(), 1)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if dbChat.ID != 1 {
		t.Errorf("Failed to store first chat in the database\nData: %+v", dbChat)
	}
	if cmd.Result.Chat.ID != 1 {
		t.Errorf("Created chat with wrong ID - %d\n", cmd.Result.Chat.ID)
	}
	if cmd.Result.Message.Content != "hello-world" {
		t.Errorf("Created message has wrong content - %q", cmd.Result.Message.Content)
	}
}

func TestCreateChatNames(t *testing.T) {
	tooLong := "toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong"

	chatCreators := []struct {
		name         string
		setup        func() (core.Command, *sql.DB)
		expectToTrim bool
	}{
		{
			name: "create chat with message",
			setup: func() (core.Command, *sql.DB) {
				c, db := test_helpers.CreateCore()
				name := tooLong
				cmd := core.NewCreateChatWithMessageCommand(c, &models.Message{
					Role:    "user",
					Content: name,
				})
				return cmd, db
			},
			expectToTrim: true,
		},
		{
			name: "create chat and completion",
			setup: func() (core.Command, *sql.DB) {
				c, db := test_helpers.CreateCore()
				name := tooLong
				cmd := core.NewCreateChatAndCompletionCommand(
					c,
					core.AssistantRole,
					name,
					"",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
				return cmd, db
			},
			expectToTrim: true,
		},
		{
			name: "with normal name",
			setup: func() (core.Command, *sql.DB) {
				c, db := test_helpers.CreateCore()
				name := "just a normal name"
				cmd := core.NewCreateChatAndCompletionCommand(
					c,
					core.AssistantRole,
					name,
					"",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
				return cmd, db
			},
			expectToTrim: false,
		},
		{
			name: "80 characters name",
			setup: func() (core.Command, *sql.DB) {
				c, db := test_helpers.CreateCore()
				name := "01234567890123456789012345678901234567890123456789012345678901234567890123456789"
				cmd := core.NewCreateChatAndCompletionCommand(
					c,
					core.AssistantRole,
					name,
					"",
					&ai_clients.Parameters{Model: "gpt-4o"},
					test_helpers.MockCompletion,
				)
				return cmd, db
			},
			expectToTrim: false,
		},
	}

	for _, test := range chatCreators {
		cmd, db := test.setup()
		err := cmd.Execute(context.Background())
		if err != nil {
			t.Errorf("Unexpected error in %q: %s\n", test.name, err)
		}
		dbChat, err := db_helpers.GetChatByID(db, context.Background(), 1)
		if err != nil {
			t.Fatalf("Unexpected error in %q: %s\n", test.name, err)
		}

		if test.expectToTrim {
			if len(dbChat.Name) > 80 {
				t.Errorf(
					"Failed to trim chat name in %q\nData: %+v\n",
					test.name,
					dbChat,
				)
			}
			if dbChat.Name[len(dbChat.Name)-3:] != "..." {
				t.Errorf(
					"Name did not ended in ellipsis in %q\nData: %+v\n",
					test.name,
					dbChat,
				)
			}
			continue
		}

		if dbChat.Name[len(dbChat.Name)-3:] == "..." {
			t.Errorf(
				"Ellipsis used on name wihtout need in %q\nData: %+v\n",
				test.name,
				dbChat,
			)
		}
	}
}
