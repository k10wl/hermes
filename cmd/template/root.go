package template

import (
	"github.com/spf13/cobra"
)

var TemplateCommand = &cobra.Command{
	Use:   "template",
	Short: "Manage custom template",
	Long:  `Manage custom template for use in chats, providing prompts or context for your interactions. Templates can be in various formats, including nested template, structs, JSON, or plain text.`,
	Example: `  $ hermes template upsert --content '--{{define "tldr"}}tldr--{{end}}'
  $ hermes template view   --name tldr
  $ hermes template edit   --name tldr
  $ hermes template delete --name tldr`,
}

func init() {
	TemplateCommand.AddCommand(deleteCommand)
	TemplateCommand.AddCommand(editCommand)
	TemplateCommand.AddCommand(upsertCommand)
	TemplateCommand.AddCommand(viewCommand)
}
