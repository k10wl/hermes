package template

import (
	"fmt"

	"github.com/k10wl/hermes/cmd/utils"
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

func createViewCommand(c *core.Core) *cobra.Command {
	viewCommand := &cobra.Command{
		Use:   "view",
		Short: "Display the contents of a specified template.",
		Long: `Retrieves template from the database and shows what is stored. Accepts an optional ` + "`--name`" + ` ` + "(`-n`)" + ` parameter - a string with SQL regex for the name of the template. If only one stored template matches the regex string, it shows the content. If the regex has multiple matches, it returns a list of matches. If the name was not provided, it returns a list of all templates.
`,
		Example: `$ hermes template view
$ hermes template view -n tldr
$ hermes template view -n %`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			query := core.NewGetTemplatesByRegexp(c, name)
			if err := query.Execute(cmd.Context()); err != nil {
				return err
			}
			config := c.GetConfig()
			if len(query.Result) == 1 {
				fmt.Fprintf(config.Stdoout, "%s", query.Result[0].Content)
				return nil
			}
			utils.ListTemplates(query.Result, config.Stdoout)
			return nil
		},
	}

	viewCommand.Flags().StringP(
		"name",
		"n",
		"%",
		"SQL regexp for name",
	)

	return viewCommand
}
