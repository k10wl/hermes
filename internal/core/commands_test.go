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

	if err := core.NewCreateTemplateCommand(
		coreInstance,
		`{{define "welcome"}}hello world!{{end}}`,
	).Execute(context.Background()); err != nil {
		panic(err)
	}
	if err := core.NewCreateTemplateCommand(
		coreInstance,
		`{{define "wrapper"}}wrapper - {{.}} - wrapper{{end}}`,
	).Execute(context.Background()); err != nil {
		panic(err)
	}
	if err := core.NewCreateTemplateCommand(
		coreInstance,
		`{{define "nested1"}}nested1: {{template "nested2" .}}{{end}}`,
	).Execute(context.Background()); err != nil {
		panic(err)
	}
	if err := core.NewCreateTemplateCommand(
		coreInstance,
		`{{define "nested2"}}nested2: {{.}}{{end}}`,
	).Execute(context.Background()); err != nil {
		panic(err)
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
	}

	for _, test := range table {
		test.init()
		err := currentCommand.Execute(context.TODO())
		test.expectedResult.TimestampsToNilForTest__()
		currentCommand.Result.TimestampsToNilForTest__()
		if test.shouldError && err == nil {
			t.Errorf("%q expected to error, but did not\n", test.name)
			continue
		}
		if !test.shouldError && err != nil {
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
	var command *core.CreateTemplateCommand

	table := []testCase{
		{
			name:         "create welcome template",
			template:     `{{define "welcome"}}hello world!{{end}}`,
			templateName: []string{"welcome"},
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
