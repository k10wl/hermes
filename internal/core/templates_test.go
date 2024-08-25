package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/db"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
)

func TestExtractTemplateDefinitionName(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected string
		errors   bool
	}
	table := []testCase{
		{
			name:     "should get template name",
			input:    `{{define "welcome"}}`,
			expected: "welcome",
			errors:   false,
		},
		{
			name: "should fail on bad string",
			input: `{{define "we
            lcome"}}`,
			expected: "",
			errors:   true,
		},
		{
			name:     "should fail on bad string",
			input:    `{{define "we`,
			expected: "",
			errors:   true,
		},
		{
			name:     "should fail on bad string",
			input:    ``,
			expected: "",
			errors:   true,
		},
		{
			name:     "should get template name if there is nested template",
			input:    `{{define "welcome"}}{{template "stuff"}}{{end}}`,
			expected: "welcome",
			errors:   false,
		},
	}
	for _, test := range table {
		actual, err := extractTemplateDefinitionName(test.input)
		if test.errors && err == nil {
			t.Errorf("%q expected error, but got nil\n", test.name)
			continue
		}
		if !test.errors && err != nil {
			t.Errorf("%q unexpected error: %v\n", test.name, err)
			continue
		}
		if actual != test.expected {
			t.Errorf(
				"%q bad result.\nexpected: %q\nactual:   %q\n",
				test.name,
				test.expected,
				actual,
			)
		}
	}
}

func TestExtractUsedTemplates(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected []string
		errors   bool
	}
	table := []testCase{
		{
			name: "should get template name",
			input: `{{template "name1"}}{{template "name2"}}

            {{template "name3"}}`,
			expected: []string{"name1", "name2", "name3"},
			errors:   false,
		},
		{
			name: "should get template name with dot or variable assigned",
			input: `{{template "name1" .Variable}}{{template "name2" .}}

            {{template "name3" .variable}}`,
			expected: []string{"name1", "name2", "name3"},
			errors:   false,
		},
		{
			name:     "return no matches on regular string",
			input:    `this is regular string`,
			expected: []string{},
			errors:   false,
		},
	}
	for _, test := range table {
		actual, err := extractTemplates(test.input)
		if test.errors && err == nil {
			t.Errorf("%q expected error, but got nil\n", test.name)
			continue
		}
		if !test.errors && err != nil {
			t.Errorf("%q unexpected error: %v\n", test.name, err)
			continue
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf(
				"%q bad result.\nexpected: %+v\nactual:   %+v\n",
				test.name,
				test.expected,
				actual,
			)
		}
	}

}

func TestBuildTemplateString(t *testing.T) {
	type testCase struct {
		name        string
		input       []string
		template    string
		expected    []string
		shouldError bool
	}

	coreInstance, _ := createCoreAndDB()

	dbTemplates := map[string]string{
		"hello-world": `{{define "hello-world"}}Hello world!{{end}}`,
		"bye-world":   `{{define "bye-world"}}Goodbye, cruel world!{{end}}`,
		"nested-l1":   `{{define "nested-l1"}}nested l1 {{template "nested-l2"}}{{end}}`,
		"nested-l2":   `{{define "nested-l2"}}nested l2 {{template "nested-l3"}}{{end}}`,
		"nested-l3":   `{{define "nested-l3"}}nested l3{{end}}`,
		"loop1":       `{{define "loop1"}}{{template "loop2"}}{{end}}`,
		"loop2":       `{{define "loop2"}}{{template "loop1"}}{{end}}`,
	}
	for _, template := range dbTemplates {
		if err := NewCreateTemplateCommand(coreInstance, template).Execute(context.Background()); err != nil {
			panic(err)
		}
	}

	table := []testCase{
		{
			name:        "should return empty string if no template was specified",
			input:       []string{"raw string"},
			template:    "",
			expected:    []string{},
			shouldError: false,
		},
		{
			name: "should return empty input if it contains unknown template",
			input: []string{
				`{{template "theduck"}}`,
				`{{template "theduck" .}}`,
				`{{template "theduck" .Value.Value}}`,
			},
			template:    "",
			expected:    []string{},
			shouldError: false,
		},
		{
			name: "should return template definition if input contains known template",
			input: []string{
				`{{template "hello-world"}}`,
				`{{template "hello-world" .}}`,
				`{{template "hello-world" .Value.Value}}`,
			},
			template:    "",
			expected:    []string{"hello-world"},
			shouldError: false,
		},
		{
			name: "should return multiple template definitions",
			input: []string{
				`{{template "hello-world"}}{{template "bye-world"}}`,
				`{{template "hello-world" .}}{{template "bye-world" .}}`,
				`{{template "hello-world" .Value.Value}}{{template "bye-world" .Value.Value}}`,
			},
			template:    "",
			expected:    []string{"hello-world", "bye-world"},
			shouldError: false,
		},
		{
			name: "should return multiple template definitions",
			input: []string{
				`{{template "hello-world"}}{{template "bye-world"}}`,
				`{{template "hello-world" .}}{{template "bye-world" .}}`,
				`{{template "hello-world" .Value.Value}}{{template "bye-world" .Value.Value}}`,
			},
			template:    "",
			expected:    []string{"hello-world", "bye-world"},
			shouldError: false,
		},
		{
			name: "should return nested template definitions",
			input: []string{
				`{{template "nested-l1"}}`,
			},
			template:    "",
			expected:    []string{"nested-l1", "nested-l2", "nested-l3"},
			shouldError: false,
		},
		{
			name: "should process looped nested templates",
			input: []string{
				`{{template "loop1"}}`,
			},
			template:    "",
			expected:    []string{"loop1", "loop2"},
			shouldError: false,
		},
		{
			name: "should return template specified in argument",
			input: []string{
				`hi`,
			},
			template:    "hello-world",
			expected:    []string{"hello-world"},
			shouldError: false,
		},
		{
			name: "should return template specified in argument with input templates",
			input: []string{
				`{{template "bye-world"}}`,
			},
			template:    "hello-world",
			expected:    []string{"hello-world", "bye-world"},
			shouldError: false,
		},
		{
			name: "should return nested template definitions with template as an argument",
			input: []string{
				`{{template "nested-l1"}}`,
			},
			template:    "hello-world",
			expected:    []string{"nested-l1", "nested-l2", "nested-l3", "hello-world"},
			shouldError: false,
		},
		{
			name: "should error if template name provided by an argument does not exist",
			input: []string{
				`{{template "nested-l1"}}`,
			},
			template:    "bullshit",
			expected:    []string{},
			shouldError: true,
		},
	}

	for _, test := range table {
		for _, input := range test.input {
			templateBuilder := newTemplateBuilder(coreInstance)
			templateBuilder.mustProcessTemplate(test.template)
			templateBuilder.process(context.Background(), input)
			processedTemplate, err := templateBuilder.string()
			if test.shouldError {
				if err == nil {
					t.Errorf("%q expected error, received nil\n", test.name)
				}
				continue
			}
			if err != nil {
				t.Errorf("%q unexpected error: %v\n", test.name, err)
				continue
			}
			for _, name := range test.expected {
				expectedTemplate, ok := dbTemplates[name]
				if !ok {
					panic(fmt.Sprintf("%q bad test setup, no template %q\n", test.name, name))
				}
				if !strings.Contains(processedTemplate, expectedTemplate) {
					t.Errorf(
						"%q bad result.\nexpected: %+v\nactual:   %+v\n",
						test.name,
						test.expected,
						processedTemplate,
					)
					break
				}
			}
		}
	}

}

type MockAIClient struct{}

func (mockClient MockAIClient) ChatCompletion(
	messages []ai_clients.Message,
) (ai_clients.Message, int, error) {
	messages[0].Role = AssistantRole
	return messages[0], 1, nil
}

func createCoreAndDB() (*Core, db.Client) {
	db, err := sqlite3.NewSQLite3(&settings.Config{DatabaseDSN: ":memory:"})
	if err != nil {
		panic(err)
	}
	return NewCore(MockAIClient{}, db), db
}
