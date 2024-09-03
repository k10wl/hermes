package main

import (
	"fmt"
	"io"
	"os"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/launch"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
	client "github.com/k10wl/openai-client"
)

var getConfig = settings.GetConfig
var newOpenAIAdapter = ai_clients.NewOpenAIAdapter

func run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	config, err := getConfig(stdin, stdout, stderr)
	if err != nil {
		return err
	}
	// FIXME we don't need fucking adapter here
	openai := newOpenAIAdapter(
		client.NewOpenAIClient(config.OpenAIKey),
	)
	// FIXME do not set model
	err = openai.SetModel(config.Model)
	if err != nil {
		return err
	}
	sqlite, err := sqlite3.NewSQLite3(config)
	if err != nil {
		return err
	}
	defer sqlite.Close()
	hermesCore := core.NewCore(openai, sqlite)
	return launch.PickStrategy(config).Execute(hermesCore, config)
}

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
