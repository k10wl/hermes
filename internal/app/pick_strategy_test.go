package app

import (
	"fmt"
	"testing"

	"github.com/k10wl/hermes/internal/runtime"
)

func TestPickStrategy(t *testing.T) {
	type tc struct {
		name     string
		input    []runtime.Config
		expected interface{}
	}

	table := []tc{
		{
			name: "return web app launcher",
			input: []runtime.Config{
				{
					Web: true,
				},
				{
					Web:    true,
					Prompt: "any string",
					Host:   runtime.DefaultHost,
					Port:   runtime.DefaultPort,
					Last:   true,
				},
			},
			expected: &launchWeb{},
		},
		{
			name: "return CLI launcher",
			input: []runtime.Config{
				{
					Prompt: "this is my prompt",
					Host:   runtime.DefaultHost,
					Port:   runtime.DefaultPort,
				},
			},
			expected: &launchCLI{},
		},
		{
			name: "return bad input launcher",
			input: []runtime.Config{
				{
					Prompt: "",
					Host:   runtime.DefaultHost,
					Port:   runtime.DefaultPort,
				},
				{
					Prompt: "     ",
					Host:   runtime.DefaultHost,
					Port:   runtime.DefaultPort,
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
