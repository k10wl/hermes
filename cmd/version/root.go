package version

import (
	"fmt"

	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

func CreateVersionCommand(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of hermes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(
				c.GetConfig().Stdoout,
				"Name:         hermes\nVersion:      %s\nVersion date: %s\n",
				c.GetConfig().Version,
				c.GetConfig().VersionDate,
			)
		},
	}
}
