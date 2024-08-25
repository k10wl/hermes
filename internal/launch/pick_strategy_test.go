package launch

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
					Web:     true,
					Content: "any string",
					Host:    settings.DefaultHostname,
					Port:    settings.DefaultPort,
					Last:    true,
				},
			},
			expected: &launchWeb{},
		},
		{
			name: "return CLI launcher",
			input: []settings.Config{
				{
					Content: "this is my prompt",
					Host:    settings.DefaultHostname,
					Port:    settings.DefaultPort,
				},
			},
			expected: &launchCLI{},
		},
		{
			name: "return bad input launcher",
			input: []settings.Config{
				{
					Content: "",
					Host:    settings.DefaultHostname,
					Port:    settings.DefaultPort,
				},
				{
					Content: "     ",
					Host:    settings.DefaultHostname,
					Port:    settings.DefaultPort,
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
					Content: "",
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
