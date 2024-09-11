package template

import (
	"context"
	"fmt"

	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

func createDeleteCommand(c *core.Core) *cobra.Command {
	deleteCommand := &cobra.Command{
		Use:     "delete",
		Short:   "Remove a template by the specified name.",
		Long:    `Mark the template with the given name as deleted. Ensure that the template you wish to delete is not currently in use. It expects the ` + "`--name`" + ` or ` + "`-n`" + ` flag to indicate which template must be deleted.`,
		Example: `$ hermes template delete -n tldr`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := c.GetConfig()
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			command := core.NewDeleteTemplateByName(c, name)
			if err := command.Execute(context.Background()); err != nil {
				return err
			}
			fmt.Fprintf(config.Stdoout, "Template %q successfully deleted.\n", name)
			return nil
		},
	}

	deleteCommand.Flags().StringP(
		"name",
		"n",
		"",
		"exact name of the template to be deleted",
	)
	err := deleteCommand.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}

	return deleteCommand
}
