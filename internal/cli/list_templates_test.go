package cli

import (
	"strings"
	"testing"

	"github.com/k10wl/hermes/internal/models"
)

func TestListTemplates(t *testing.T) {
	type testCase struct {
		name     string
		input    []*models.Template
		expected string
	}

	table := []testCase{
		{
			name: "should return list of templates",
			input: []*models.Template{
				{
					Name:    "first",
					Content: `--{{define "first"}}first--{{end}}`,
				},
				{
					Name:    "second",
					Content: `--{{define "second"}}second--{{end}}`,
				},
			},
			expected: `List of templates:
[Name]    first
[Content] --{{define "first"}}first--{{end}}
--------------------
[Name]    second
[Content] --{{define "second"}}second--{{end}}
--------------------
`,
		},
		{
			name:     "should return no templates if none are stored",
			input:    []*models.Template{},
			expected: "No templates. Please use -h to get info of how to add templates\n",
		},
	}

	sb := &strings.Builder{}
	for _, test := range table {
		sb.Reset()
		err := listTemplates(test.input, sb)
		if err != nil {
			t.Errorf("%q - unexpected error: %v\n\n", test.name, err)
			continue
		}
		actual := sb.String()
		if actual != test.expected {
			t.Errorf(
				"%q - bad result\nexpected: %q\nactual:   %q\n\n",
				test.name,
				test.expected,
				actual,
			)
		}
	}
}
