package templates

import (
	"fmt"

	"github.com/k10wl/hermes/cmd/utils"
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

var viewCommand = &cobra.Command{
	Use:   "view",
	Short: "Display the contents of a specified template.",
	Long: `Retreives templates from database and shows what is stored.
Accepts optional --name -n param - string with SQL regexp for name of template.
If only one stored template matches regexp string - shows content.
If regexp has multiple matches - returns list of matches.
If name was not provided - returns list all templates.`,
	Example: `$ hermes templates view
$ hermes templates view -n tldr
$ hermes templates view -n %`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			panic(err)
		}
		c := utils.GetCore(cmd)
		query := core.NewGetTemplatesByRegexp(c, name)
		if err := query.Execute(cmd.Context()); err != nil {
			panic(err)
		}
		if len(query.Result) == 1 {
			fmt.Fprintf(c.GetConfig().Stdoout, "%s", query.Result[0].Content)
			return
		}
		utils.ListTemplates(query.Result, c.GetConfig().Stdoout)
	},
}

func init() {
	viewCommand.Flags().StringP(
		"name",
		"n",
		"%",
		"SQL regexp for name",
	)
}
