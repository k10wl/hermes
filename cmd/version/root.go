package version

import (
	"fmt"

	"github.com/k10wl/hermes/cmd/utils"
	"github.com/spf13/cobra"
)

var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of hermes",
	Long:  "All software has versions. This is Hugo's",
	Run: func(cmd *cobra.Command, args []string) {
		c := utils.GetCore(cmd)
		fmt.Fprintf(
			c.GetConfig().Stdoout,
			"Name:         hermes\nVersion:      %s\nVersion date: %s\n",
			c.GetConfig().Version,
			c.GetConfig().VersionDate,
		)
	},
}
