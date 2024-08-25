package launch

import (
	"testing"

	"github.com/k10wl/hermes/internal/settings"
)

func TestPickStrategy(t *testing.T) {
	type tc struct {
		name     string
		input    []settings.Config
		expected interface{}
	}

	table := []tc{
		{
			name: "return web app launcher",
			input: []settings.Config{
				{
					WebFlags: settings.WebFlags{
						Web: true,
					},
				},
				{
					WebFlags: settings.WebFlags{
						Web:  true,
						Host: settings.DefaultHostname,
						Port: settings.DefaultPort,
					},
					CLIFlags: settings.CLIFlags{
						Content: "any string",
						Last:    true,
					},
				},
			},
			expected: &launchWeb{},
		},
		{
			name: "return CLI launcher",
			input: []settings.Config{
				{
					CLIFlags: settings.CLIFlags{
						Content: "this is my prompt",
					},
					WebFlags: settings.WebFlags{
						Host: settings.DefaultHostname,
						Port: settings.DefaultPort,
					},
				},
				{
					TemplateFlags: settings.TemplateFlags{
						UpsertTemplate: `{{define "hi"}}hello world{{end}}`,
					},
				},
			},
			expected: &launchCLI{},
		},
		{
			name: "return bad input launcher",
			input: []settings.Config{
				{
					CLIFlags: settings.CLIFlags{
						Content: "",
					},
					WebFlags: settings.WebFlags{
						Host: settings.DefaultHostname,
						Port: settings.DefaultPort,
					},
				},
				{
					CLIFlags: settings.CLIFlags{
						Content: "     ",
					},
					WebFlags: settings.WebFlags{
						Host: settings.DefaultHostname,
						Port: settings.DefaultPort,
					},
				},
				{
					WebFlags: settings.WebFlags{
						Web: false,
					},
				},
				{
					WebFlags: settings.WebFlags{
						Web: false,
					},
				},
				{
					CLIFlags: settings.CLIFlags{
						Last: true,
					},
					WebFlags: settings.WebFlags{
						Web: false,
					},
				},
				{
					CLIFlags: settings.CLIFlags{
						Content: "",
					},
				},
			},
			expected: &launchBadInput{},
		},
	}

	for _, test := range table {
		for _, input := range test.input {
			output := PickStrategy(&input)
			if output != test.expected {
				t.Errorf(
					"%s:\nexpected: %T\nactual:   %T",
					test.name,
					test.expected,
					output,
				)
			}
		}
	}
}
