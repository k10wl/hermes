package templates

import (
	"context"
	"fmt"

	"github.com/k10wl/hermes/cmd/utils"
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Remove a template by the specified name.",
	Long: `Mark template with given name as deleted.
Ensure that the template you wish to delete is not currently in use.
Expects --name -n flag to indicate what template must be deleted.`,
	Example: `$ hermes templates delete -n tldr
Template "test2" successfully deleted.

$ hermes templates delete -n tldr
Failed. Template "test2" not found. `,
	Run: func(cmd *cobra.Command, args []string) {
		c := utils.GetCore(cmd)
		config := c.GetConfig()
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			utils.LogError(config.Stderr, err)
			return
		}
		command := core.NewDeleteTemplateByName(c, name)
		if err := command.Execute(context.Background()); err != nil {
			utils.LogError(config.Stderr, err)
			return
		}
		fmt.Fprintf(config.Stdoout, "Template %q successfully deleted.\n", name)
	},
}

func init() {
	deleteCommand.Flags().StringP(
		"name",
		"n",
		"",
		"exact name of template to be deleted **required**",
	)
	err := deleteCommand.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}
}
