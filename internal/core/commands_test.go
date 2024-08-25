package core_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
)

func TestCreateChatAndCompletionCommand(t *testing.T) {
	type testCase struct {
		name           string
		init           func()
		shouldError    bool
		expectedResult models.Message
	}

	coreInstance, _ := createCoreAndDB()
	var currentCommand *core.CreateChatAndCompletionCommand

	dbTemplates := map[string]string{
		"welcome": `{{define "welcome"}}hello world!{{end}}`,
		"wrapper": `{{define "wrapper"}}wrapper - {{.}} - wrapper{{end}}`,
		"nested1": `{{define "nested1"}}nested1: {{template "nested2" .}}{{end}}`,
		"nested2": `{{define "nested2"}}nested2: {{.}}{{end}}`,
		"loop1":   `{{define "loop1"}}{{template "loop2"}}{{end}}`,
		"loop2":   `{{define "loop2"}}{{template "loop1"}}{{end}}`,
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
		{
			name: "Should fill template data when template is passed as argument",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance, core.AssistantRole, ``, "welcome",
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      4,
				ChatID:  2,
				Role:    core.AssistantRole,
				Content: "hello world!",
			},
		},
		{
			name: "Should fill inner template data when template name is not provided",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance, core.AssistantRole, `{{template "welcome"}}`, "",
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      6,
				ChatID:  3,
				Role:    core.AssistantRole,
				Content: "hello world!",
			},
		},
		{
			name: "Should fill inner template data and provided template name data",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance, core.AssistantRole, `{{template "welcome"}}`, "wrapper",
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      8,
				ChatID:  4,
				Role:    core.AssistantRole,
				Content: "wrapper - hello world! - wrapper",
			},
		},
		{
			name: "Should error if given template name does not exist",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance, core.AssistantRole, `{{template "welcome"}}`, "does not exist",
				)
			},
			shouldError:    true,
			expectedResult: models.Message{},
		},
		{
			name: "Should process string input and fill template data",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance, core.AssistantRole, "hello world!", "wrapper",
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ID:      10,
				ChatID:  5,
				Role:    core.AssistantRole,
				Content: "wrapper - hello world! - wrapper",
			},
		},
		{
			name: "Should error on circular templates with provided template",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance, core.AssistantRole, "will blow up", "loop1",
				)
			},
			shouldError:    true,
			expectedResult: models.Message{},
		},
		{
			name: "Should error on circular templates with inputted template",
			init: func() {
				currentCommand = core.NewCreateChatAndCompletionCommand(
					coreInstance, core.AssistantRole, `will blow up {{template "loop1" . }} `, "",
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
					`will blow up {{template "loop1" . }} `,
					"loop2",
				)
			},
			shouldError:    true,
			expectedResult: models.Message{},
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
				t.Errorf("%q expected to error, but did not\n", test.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q unexpected error: %v\n", test.name, err)
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

	coreInstance, _ := createCoreAndDB()
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
				ChatID:  1,
				Role:    core.AssistantRole,
				Content: "hello world",
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

	coreInstance, db := createCoreAndDB()
	var command *core.UpsertTemplateCommand

	table := []testCase{
		{
			name:         "create welcome template",
			template:     `{{define "welcome"}}hello world!{{end}}`,
			templateName: []string{"welcome"},
			init: func() {
				command = core.NewUpsertTemplateCommand(
					coreInstance,
					`{{define "welcome"}}hello world!{{end}}`,
				)
			},
			shouldError: false,
		},
		{
			name:         "should override written command",
			template:     `{{define "welcome"}}welcome world!{{end}}`,
			templateName: []string{"welcome"},
			init: func() {
				command = core.NewUpsertTemplateCommand(
					coreInstance,
					`{{define "welcome"}}welcome world!{{end}}`,
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
