package template

import (
	"context"
	"fmt"
	"os"

	"github.com/k10wl/hermes/cmd/utils"
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

var editCommand = &cobra.Command{
	Use:   "edit",
	Short: "Modify an existing template by its name",
	Long: `Allows you to edit the stored template. This command finds a template with the name provided in the flag --name (-n) and opens your preferred editor with its content. Upon saving and closing the file, the edit will be stored in the database.

Behavior:
1. If the edited content is identical to the original, no changes are made.
2. If the edited template's name differs:
    - If ` + "`--clone`" + ` is true and the new name is unique, a new template is created, retaining the original.
    - If ` + "`--clone`" + ` is false, and the new name is unique, the template gets renamed.
    - If a name conflict arises, an error will be returned.
3. If the edited content is invalid, an error will be returned.
`,
	Example: `$ hermes template edit --name tldr
$ hermes template edit --name tldr --clone`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		c := utils.GetCore(cmd)
		config := c.GetConfig()
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		clone, err := cmd.Flags().GetBool("clone")
		if err != nil {
			return err
		}
		query := core.NewGetTemplatesByNamesQuery(c, []string{name})
		if err := query.Execute(ctx); err != nil {
			return err
		}
		if len(query.Result) != 1 {
			return fmt.Errorf("template %q not found\n", name)
		}
		editedContent, err := utils.OpenInEditor(
			query.Result[0].Content,
			os.Stdin,
			config.Stdoout,
			config.Stderr,
		)
		if err != nil {
			return err
		}
		if editedContent == query.Result[0].Content {
			return fmt.Errorf("edit is identical to original\n")
		}
		if err := core.NewEditTemplateByName(
			c,
			name,
			editedContent,
			clone,
		).Execute(ctx); err != nil {
			return err
		}
		if clone {
			fmt.Fprintf(config.Stdoout, "Template cloned and edited successfully\n")
			return nil
		}
		fmt.Fprintf(config.Stdoout, "Template edited successfully\n")
		return nil
	},
}

func init() {
	editCommand.Flags().StringP(
		"name",
		"n",
		"",
		"exact name of the template to be edited",
	)
	err := editCommand.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}
	editCommand.Flags().BoolP(
		"clone",
		"c",
		false,
		"keep the original template; returns an error if the name can't be updated",
	)
}
