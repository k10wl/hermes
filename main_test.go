package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/cli"
	"github.com/k10wl/hermes/internal/settings"
	client "github.com/k10wl/openai-client"
)

var completion = "This is a mock response."

func TestApp(t *testing.T) {
	type expected struct {
		stdout  string
		stderr  string
		network string
	}
	type tc struct {
		name        string
		expected    expected
		shouldError bool
		prepare     func()
		extraFn     func() error
	}

	oldConfig := getConfig
	oldOpenAIAdapter := newOpenAIAdapter
	newOpenAIAdapter = func(client *client.OpenAIClient) ai_clients.OpenAIAdapterInterface {
		return &MockOpenAIAdapter{}
	}
	defer func() {
		getConfig = oldConfig
		newOpenAIAdapter = oldOpenAIAdapter
	}()
	var stdin bytes.Buffer
	var stdout strings.Builder
	var stderr strings.Builder
	reset := func() {
		stdin.Reset()
		stdout.Reset()
		stderr.Reset()
	}

	table := []tc{
		{
			name: "output help info if no input given",
			prepare: func() {
				getConfig = func(
					stdin io.Reader,
					stdout io.Writer,
					stderr io.Writer,
				) (*settings.Config, error) {
					c, err := oldConfig(stdin, stdout, stderr)
					c.DatabaseDSN = ":memory:"
					return c, err
				}
			},
			expected: expected{stdout: cli.GetHelpString(settings.Version) + "\n"},
		},

		{
			name: "complete message",
			prepare: func() {
				getConfig = func(
					stdin io.Reader,
					stdout io.Writer,
					stderr io.Writer,
				) (*settings.Config, error) {
					c, err := oldConfig(stdin, stdout, stderr)
					c.DatabaseDSN = ":memory:"
					c.Content = "complete prompt"
					return c, err
				}
			},
			expected: expected{stdout: completion + "\n"},
		},

		{
			name: "start web server",
			prepare: func() {
				getConfig = func(
					stdin io.Reader,
					stdout io.Writer,
					stderr io.Writer,
				) (*settings.Config, error) {
					c, err := oldConfig(stdin, stdout, stderr)
					c.Port = "8124"
					c.DatabaseDSN = ":memory:"
					c.Content = "complete prompt"
					c.Web = true
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					c.ShutdownContext = ctx
					return c, err
				}
			},
			expected: expected{
				stdout: "Starting server on 127.0.0.1:8124\nShutdown signal received\n",
			},
			extraFn: func() error {
				resp, err := http.Get("http://127.0.0.1:8124")
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				return nil
			},
		},
	}

	for _, test := range table {
		test.prepare()
		err := run(&stdin, &stdout, &stderr)
		if test.shouldError && err == nil {
			t.Errorf("Completed, but expected error: %s", test.name)
		}
		if !test.shouldError && err != nil {
			t.Errorf("Unexpected error in: %s\n Error: %v", test.name, err)
		}
		if test.expected.stdout != stdout.String() {
			t.Errorf(
				"Failed stdout: %s\n---\nexpected: %q\n---\nactual:   %q\n",
				test.name,
				test.expected.stdout,
				stdout.String(),
			)
		}
		if test.expected.stderr != stderr.String() {
			t.Errorf(
				"Failed stderr: %s\n---\nexpected: %q\n---\nactual:   %q\n",
				test.name,
				test.expected.stdout,
				stdout.String(),
			)
		}
		if test.extraFn != nil {
			err := test.extraFn()
			if err != nil {
				// XXX damn i fucked up, this test requires active web server
				t.Errorf("Failed extraFn test for %s\nError: %v", test.name, err)
			}
		}
		reset()
	}
}

type MockOpenAIAdapter struct{}

func (m *MockOpenAIAdapter) ChatCompletion(
	messages []ai_clients.Message,
) (ai_clients.Message, int, error) {
	return ai_clients.Message{
			Role:    "assistant",
			Content: completion,
		}, len(
			messages,
		), nil
}

func (m *MockOpenAIAdapter) SetModel(model string) error {
	return nil
}
