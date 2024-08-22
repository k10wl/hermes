package app

import (
	"fmt"
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
					Web: true,
				},
				{
					Web:    true,
					Prompt: "any string",
					Host:   settings.DefaultHost,
					Port:   settings.DefaultPort,
					Last:   true,
				},
			},
			expected: &launchWeb{},
		},
		{
			name: "return CLI launcher",
			input: []settings.Config{
				{
					Prompt: "this is my prompt",
					Host:   settings.DefaultHost,
					Port:   settings.DefaultPort,
				},
			},
			expected: &launchCLI{},
		},
		{
			name: "return bad input launcher",
			input: []settings.Config{
				{
					Prompt: "",
					Host:   settings.DefaultHost,
					Port:   settings.DefaultPort,
				},
				{
					Prompt: "     ",
					Host:   settings.DefaultHost,
					Port:   settings.DefaultPort,
				},
				{
					Web: false,
				},
				{
					Web: false,
				},
				{
					Web:  false,
					Last: true,
				},
				{
					Prompt: "",
				},
			},
			expected: &launchBadInput{},
		},
	}

	for _, test := range table {
		for _, input := range test.input {
			output := PickStrategy(&input)
			fmt.Printf("output: %v\n", output)
			fmt.Printf("test.expected: %v\n", test.expected)
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
