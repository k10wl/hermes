package template

import (
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

func CreateTemplateCommand(c *core.Core) *cobra.Command {
	templateCommand := &cobra.Command{
		Use:   "template",
		Short: "Manage custom templates",
		Long:  `Manage custom templates for use in chats, providing prompts or context for your interactions. Templates can be in various formats, including nested templates, structs, JSON, or plain text.`,
		Example: `  $ hermes template upsert --content '--{{define "tldr"}}tldr--{{end}}'
  $ hermes template view   --name tldr
  $ hermes template edit   --name tldr
  $ hermes template delete --name tldr`,
	}

	templateCommand.AddCommand(createDeleteCommand(c))
	templateCommand.AddCommand(createEditCommand(c))
	templateCommand.AddCommand(createUpsertCommand(c))
	templateCommand.AddCommand(createViewCommand(c))

	return templateCommand
}
