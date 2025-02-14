package template

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/k10wl/hermes/cmd/utils"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
	"github.com/spf13/cobra"
)

func relayUpsert(c *core.Core, tmp *models.Template) {
	uid := uuid.NewString()
	if data, err := messages.Encode(
		messages.NewServerTemplateCreated(uid, tmp),
	); err == nil {
		utils.NotifyActiveSessions(c, uid, data)
	}
}

func createUpsertCommand(c *core.Core) *cobra.Command {
	upsertCommand := &cobra.Command{
		Use:   `upsert`,
		Short: "Update an existing template or create a new one if it does not exist.",
		Long: `Receives template, parses, verifies, and saves content into the database. If the ` + "`--content`" + ` (` + "`-c`" + `) flag is not provided, the default text editor will be opened. The name of the template will derive from the definition/block name. The template must comply with Golang text template rules. NOTE: Delimiters differ from Golang text template:
    - Left delimiter - '--{{';
    - Right delimiter - '}}';`,
		Example: `$ hermes template upsert
$ hermes template upsert -c "--{{define "template"}}(instruction)--{{end}}
$ hermes template upsert --content "--{{define "template"}}(instruction)--{{end}} `,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := c.GetConfig()
			content, err := cmd.Flags().GetString("content")
			if err != nil {
				return err
			}
			if content == "" {
				content, err = utils.OpenInEditor(
					models.DefaultTemplate.Content,
					config.Stdin,
					config.Stdoout,
					config.Stderr,
				)
				if content == models.DefaultTemplate.Content {
					return fmt.Errorf("do not save example template, make some changes\n")
				}
				if err != nil {
					return err
				}
			}
			upsertCmd := core.NewUpsertTemplateCommand(
				c,
				content,
			)
			if err := upsertCmd.Execute(context.Background()); err != nil {
				return err
			}
			relayUpsert(c, upsertCmd.Result)
			fmt.Fprintf(config.Stdoout, "Template upserted successfully\n")
			return nil
		},
	}

	upsertCommand.Flags().StringP(
		"content",
		"c",
		"",
		"template content; if not provided, the text editor will be opened",
	)

	return upsertCommand
}
