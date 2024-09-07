package templates

import (
	"github.com/spf13/cobra"
)

var TemplatesCommand = &cobra.Command{
	Use:   "templates",
	Short: "Manage custom templates",
	Long:  `Manage custom templates for use in chats, providing prompts or context for your interactions. Templates can be in various formats, including nested templates, structs, JSON, or plain text.`,
	Example: `  $ hermes templates upsert --content '--{{define "tldr"}}tldr--{{end}}'
  $ hermes templates view   --name tldr
  $ hermes templates edit   --name tldr
  $ hermes templates delete --name tldr`,
}

func init() {
	TemplatesCommand.AddCommand(deleteCommand)
	TemplatesCommand.AddCommand(editCommand)
	TemplatesCommand.AddCommand(upsertCommand)
	TemplatesCommand.AddCommand(viewCommand)
}
