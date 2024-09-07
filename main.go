package main

import (
	"fmt"
	"io"
	"os"

	"github.com/k10wl/hermes/cmd"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
)

var getConfig = settings.GetConfig

func run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (*core.Core, error) {
	config, err := getConfig(stdin, stdout, stderr)
	if err != nil {
		return nil, err
	}
	sqlite, err := sqlite3.NewSQLite3(config)
	if err != nil {
		return nil, err
	}
	hermesCore := core.NewCore(sqlite, config)
	return hermesCore, nil
}

func main() {
	core, err := run(os.Stdin, os.Stdout, os.Stderr)
	defer core.GetDB().Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	if err := cmd.Execute(core); err != nil {
		os.Exit(1)
	}
}
