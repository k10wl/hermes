package main

import (
	"fmt"
	"io"
	"os"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/launch"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
)

var getConfig = settings.GetConfig

func run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	config, err := getConfig(stdin, stdout, stderr)
	if err != nil {
		return err
	}
	sqlite, err := sqlite3.NewSQLite3(config)
	if err != nil {
		return err
	}
	defer sqlite.Close()
	hermesCore := core.NewCore(sqlite, config)
	return launch.PickStrategy(config).Execute(hermesCore, config)
}

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
