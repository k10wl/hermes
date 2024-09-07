package utils

import (
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

func GetCore(ctx *cobra.Command) *core.Core {
	core, ok := ctx.Root().Context().Value("core").(*core.Core)
	if core == nil || !ok {
		panic("failed to get core")
	}
	return core
}
