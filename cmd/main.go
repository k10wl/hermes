package main

import (
	"fmt"
	"os"

	"github.com/k10wl/hermes/internal/app"
)

func main() {
	if err := app.Launch(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
