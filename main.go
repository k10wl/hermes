package main

import (
	"fmt"
	"io"
	"os"

	"github.com/k10wl/hermes/cmd"
	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
)

func prepare(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (*core.Core, error) {
	config, err := settings.GetConfig(stdin, stdout, stderr)
	if err != nil {
		return nil, err
	}
	sqlite, err := sqlite3.NewSQLite3(config.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	hermesCore := core.NewCore(sqlite, config)
	return hermesCore, nil
}

func main() {
	core, err := prepare(os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer core.GetDB().Close()
	if err := cmd.Execute(core, ai_clients.Complete); err != nil {
		os.Exit(1)
	}
}
