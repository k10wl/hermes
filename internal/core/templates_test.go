package core

import (
	"reflect"
	"testing"
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
