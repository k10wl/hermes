package cli

import (
	"fmt"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/runtime"
)

func CLI(core core.Core) {
	c := runtime.GetConfig()
	res, err := core.SendMessage(*c.Prompt)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
