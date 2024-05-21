package cli

import (
	"fmt"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/runtime"
)

func CLI(core *core.Core, config *runtime.Config) {
	res, err := core.SendMessage(*config.Prompt)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
