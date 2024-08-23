package core_test

import (
	"context"
	"reflect"
	"testing"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
)

func TestCreateChatAndCompletionCommand(t *testing.T) {
	type testCase struct {
		name           string
		init           func()
		shouldError    bool
		expectedResult models.Message
	}

	db, err := sqlite3.NewSQLite3(&settings.Config{DatabaseDSN: ":memory:"})
	if err != nil {
		panic(err)
	}
	coreInstance := core.NewCore(
		MockAIClient{},
		db,
	)

	var currentCommand *core.CreateChatAndCompletionCommand

	table := []testCase{
		{
			name: "Should create simple response",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance, core.AssistantRole, "hello world", "",
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      2,
				ChatID:  1,
				Role:    core.AssistantRole,
				Content: "hello world",
			},
		},
		/* TODO do this shit
		{
			name: "Should fill template data",
			init: func() {
				core.NewCreateTemplateCommand(
					coreInstance,
					`{{define "welcome"}}hello world!{{end}}`,
				).Execute(context.TODO())
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance, core.AssistantRole, `{{template "welcome"}}`, "",
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      4,
				ChatID:  2,
				Role:    core.AssistantRole,
				Content: "hello world",
			},
		},
		*/
	}

	for _, test := range table {
		test.init()
		err := currentCommand.Execute(context.TODO())
		resetMessageTime(&test.expectedResult)
		resetMessageTime(currentCommand.Result)
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

func TestCreateCompletionCommand(t *testing.T) {
	type testCase struct {
		name           string
		init           func()
		shouldError    bool
		expectedResult models.Message
	}

	db, err := sqlite3.NewSQLite3(&settings.Config{DatabaseDSN: ":memory:"})
	if err != nil {
		panic(err)
	}
	coreInstance := core.NewCore(
		MockAIClient{},
		db,
	)

	var currentCommand *core.CreateCompletionCommand

	table := []testCase{
		{
			name: "Should create simple response",
			init: func() {
				currentCommand = core.NewCreateCompletionCommand(
					coreInstance, 1, core.AssistantRole, "hello world", "",
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      2,
				Role:    core.AssistantRole,
				Content: "hello world",
			},
		},
	}

	for _, test := range table {
		test.init()
		err := currentCommand.Execute(context.TODO())
		resetMessageTime(&test.expectedResult)
		resetMessageTime(currentCommand.Result)
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
		templateName string
		init         func()
		shouldError  bool
	}
	db, err := sqlite3.NewSQLite3(&settings.Config{DatabaseDSN: ":memory:"})
	if err != nil {
		panic(err)
	}
	coreInstance := core.NewCore(
		MockAIClient{},
		db,
	)
	var command *core.CreateTemplateCommand
	table := []testCase{
		{
			name:         "create welcome template",
			template:     `{{define "welcome"}}hello world!{{end}}`,
			templateName: "welcome",
			init: func() {
				command = core.NewCreateTemplateCommand(
					coreInstance,
					`{{define "welcome"}}hello world!{{end}}`,
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
		tmpl, err := db.GetTemplateByName(context.TODO(), test.templateName)
		if err != nil {
			t.Errorf("%q failed to get created template: %v\n", test.name, err)
			continue
		}
		if test.template != tmpl.Content {
			t.Errorf(
				"%q - bad result\nexpected: %+v\nactual:   %+v\n",
				test.name,
				test.template,
				tmpl,
			)
		}
	}
}

type MockAIClient struct{}

func (mockClient MockAIClient) ChatCompletion(
	messages []ai_clients.Message,
) (ai_clients.Message, int, error) {
	messages[0].Role = core.AssistantRole
	return messages[0], 1, nil
}

func resetMessageTime(message *models.Message) {
	message.CreatedAt = nil
	message.UpdatedAt = nil
	message.DeletedAt = nil
}