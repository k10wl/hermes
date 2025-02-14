package template

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/k10wl/hermes/cmd/utils"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
	"github.com/spf13/cobra"
)

func relayEdit(c *core.Core, tmp *models.Template, action string) {
	uid := uuid.NewString()
	switch action {
	case "edit":
		if data, err := messages.Encode(
			messages.NewServerTemplateChanged(uid, tmp),
		); err == nil {
			utils.NotifyActiveSessions(c, uid, data)
		}
		break
	case "clone":
		if data, err := messages.Encode(
			messages.NewServerTemplateCreated(uid, tmp),
		); err == nil {
			utils.NotifyActiveSessions(c, uid, data)
		}
		break
	default:
		panic("unhandled relay action during edit command")
	}
}

func createEditCommand(c *core.Core) *cobra.Command {
	editCommand := &cobra.Command{
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
			editedContent, err := cmd.Flags().GetString("content")
			if err != nil || strings.TrimSpace(editedContent) == "" {
				editorOutput, err := utils.OpenInEditor(
					query.Result[0].Content,
					os.Stdin,
					config.Stdoout,
					config.Stderr,
				)
				if err != nil {
					return err
				}
				editedContent = editorOutput
			}
			if editedContent == query.Result[0].Content {
				return fmt.Errorf("edit is identical to original\n")
			}
			editCmd := core.NewEditTemplateByName(
				c,
				name,
				editedContent,
				clone,
			)
			if err := editCmd.Execute(ctx); err != nil {
				return err
			}
			if clone {
				relayEdit(c, editCmd.Result, "clone")
				fmt.Fprintf(config.Stdoout, "Template cloned and edited successfully\n")
				return nil
			}
			relayEdit(c, editCmd.Result, "edit")
			fmt.Fprintf(config.Stdoout, "Template edited successfully\n")
			return nil
		},
	}

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
	editCommand.Flags().StringP(
		"content",
		"c",
		"",
		"keep the original template; returns an error if the name can't be updated",
	)
	editCommand.Flags().Bool(
		"clone",
		false,
		"keep the original template; returns an error if the name can't be updated",
	)

	return editCommand
}
