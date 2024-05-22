package cli

import (
	"fmt"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/runtime"
)

func CLI(core *core.Core, config *runtime.Config) {
	// FIXME check if prompt is empty and print out help info if so
	res, err := core.SendMessage(config.Prompt)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
