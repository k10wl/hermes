package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestBuildTemplateString(t *testing.T) {
	type testCase struct {
		name        string
		input       []string
		template    string
		expected    []string
		shouldError bool
	}

	coreInstance, _ := __createCoreAndDB()

	dbTemplates := map[string]string{
		"hello-world": `--{{define "hello-world"}}Hello world!--{{end}}`,
		"bye-world":   `--{{define "bye-world"}}Goodbye, cruel world!--{{end}}`,
		"nested-l1":   `--{{define "nested-l1"}}nested l1 --{{template "nested-l2"}}--{{end}}`,
		"nested-l2":   `--{{define "nested-l2"}}nested l2 --{{template "nested-l3"}}--{{end}}`,
		"nested-l3":   `--{{define "nested-l3"}}nested l3--{{end}}`,
		"loop1":       `--{{define "loop1"}}--{{template "loop2"}}--{{end}}`,
		"loop2":       `--{{define "loop2"}}--{{template "loop1"}}--{{end}}`,
	}
	for _, template := range dbTemplates {
		if err := NewUpsertTemplateCommand(coreInstance, template).Execute(context.Background()); err != nil {
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
				`--{{template "theduck"}}`,
				`--{{template "theduck" .}}`,
				`--{{template "theduck" .Value.Value}}`,
			},
			template:    "",
			expected:    []string{},
			shouldError: false,
		},
		{
			name: "should return template definition if input contains known template",
			input: []string{
				`--{{template "hello-world"}}`,
				`--{{template "hello-world" .}}`,
				`--{{template "hello-world" .Value.Value}}`,
			},
			template:    "",
			expected:    []string{"hello-world"},
			shouldError: false,
		},
		{
			name: "should return multiple template definitions",
			input: []string{
				`--{{template "hello-world"}}--{{template "bye-world"}}`,
				`--{{template "hello-world" .}}--{{template "bye-world" .}}`,
				`--{{template "hello-world" .Value.Value}}--{{template "bye-world" .Value.Value}}`,
			},
			template:    "",
			expected:    []string{"hello-world", "bye-world"},
			shouldError: false,
		},
		{
			name: "should return multiple template definitions",
			input: []string{
				`--{{template "hello-world"}}--{{template "bye-world"}}`,
				`--{{template "hello-world" .}}--{{template "bye-world" .}}`,
				`--{{template "hello-world" .Value.Value}}--{{template "bye-world" .Value.Value}}`,
			},
			template:    "",
			expected:    []string{"hello-world", "bye-world"},
			shouldError: false,
		},
		{
			name: "should return nested template definitions",
			input: []string{
				`--{{template "nested-l1"}}`,
			},
			template:    "",
			expected:    []string{"nested-l1", "nested-l2", "nested-l3"},
			shouldError: false,
		},
		{
			name: "should process looped nested templates",
			input: []string{
				`--{{template "loop1"}}`,
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
				`--{{template "bye-world"}}`,
			},
			template:    "hello-world",
			expected:    []string{"hello-world", "bye-world"},
			shouldError: false,
		},
		{
			name: "should return nested template definitions with template as an argument",
			input: []string{
				`--{{template "nested-l1"}}`,
			},
			template:    "hello-world",
			expected:    []string{"nested-l1", "nested-l2", "nested-l3", "hello-world"},
			shouldError: false,
		},
		{
			name: "should error if template name provided by an argument does not exist",
			input: []string{
				`--{{template "nested-l1"}}`,
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
