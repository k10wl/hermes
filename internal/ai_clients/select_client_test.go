package ai_clients

import (
	"reflect"
	"testing"

	"github.com/k10wl/hermes/internal/settings"
)

func TestSelectProvider(t *testing.T) {
	type testCase struct {
		name        string
		input       []string
		expected    any
		shouldError bool
	}

	table := []testCase{
		{
			name: "should return gpt handler",
			input: []string{
				"openai/gpt-4.5-turbo",
				"openai/gpt-4o-mini",
				"openai/gpt-3.5-turbo",
				"openai/gpt-some-random-name",
			},
			expected: &clientOpenAI{},
		},
		{
			name: "should return claude handler",
			input: []string{
				"anthropic/claude-3-5-sonnet-20240620",
				"anthropic/claude-3-opus-20240229",
				"anthropic/claude-3-haiku-20240307",
				"anthropic/claude-some-random-name",
			},
			expected: &clientClaude{},
		},
		{
			name: "should error on unhandled provider",
			input: []string{
				"some-unhandled-provider",
				"randomstring",
			},
			expected:    nil,
			shouldError: true,
		},
	}

	for _, test := range table {
		for _, input := range test.input {
			provider, _, _ := extractProviderAndModel(input)
			res, err := selectClient(provider, &settings.Providers{
				OpenAIKey: "SECRET",
			})
			if test.shouldError {
				if err == nil {
					t.Errorf("%q - expected to error, but did not\n\n", test.name)
				}
				continue
			}
			if err != nil {
				t.Errorf("%q - unexpected error - %v", test.name, err)
				continue
			}
			expected := reflect.TypeOf(test.expected)
			actual := reflect.TypeOf(res)
			if expected != actual {
				t.Errorf(
					"%q - bad output\nexpected: %v\nactual:   %v",
					test.name,
					expected,
					actual,
				)
			}
		}
	}
}
