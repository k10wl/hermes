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
					coreInstance, core.AssistantRole, `--{{template "welcome"}}`, "",
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
					coreInstance, core.AssistantRole, `--{{template "welcome"}}`, "wrapper",
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
					coreInstance, core.AssistantRole, `--{{template "welcome"}}`, "does not exist",
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
					coreInstance, core.AssistantRole, `will blow up --{{template "loop1" . }} `, "",
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
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ChatID:  6,
				ID:      12,
				Role:    core.AssistantRole,
				Content: "should fill welcome (hello world!)(hello world!)",
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
				)
			},
			shouldError: false,
			expectedResult: models.Message{
				ChatID:  7,
				ID:      14,
				Role:    core.AssistantRole,
				Content: "wrapper - should fill welcome (hello world!)(hello world!) - wrapper",
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

	coreInstance, _ := createCoreAndDB()
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
		name         string
		init         func() string
		templateName string
		shouldError  bool
	}

	coreInstance, _ := createCoreAndDB()
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
				command = *core.NewEditTemplateByName(coreInstance, "welcome", content)
				return content
			},
			templateName: "welcome",
		},
		{
			name: "should allow block definition",
			init: func() string {
				content := `--{{block "welcome" .}}hi--{{end}}`
				command = *core.NewEditTemplateByName(coreInstance, "welcome", content)
				return content
			},
			templateName: "welcome",
		},
		{
			name: "should error if template does not exist",
			init: func() string {
				content := `--{{define "welcome"}}hi--{{end}}`
				command = *core.NewEditTemplateByName(coreInstance, "does not exist", content)
				return content
			},
			templateName: "does not exist",
			shouldError:  true,
		},
		{
			name: "should error if template has corrupted content",
			init: func() string {
				content := `--{{define "welco`
				command = *core.NewEditTemplateByName(coreInstance, "welcome", content)
				return content
			},
			templateName: "welcome",
			shouldError:  true,
		},
		{
			name: "should error if new name does not match existing name",
			init: func() string {
				content := `--{{define "missmatch"}}stuff--{{end}}`
				command = *core.NewEditTemplateByName(coreInstance, "welcome", content)
				return content
			},
			templateName: "missmatch",
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
	}
}
