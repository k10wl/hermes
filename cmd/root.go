package cmd

import (
	"github.com/k10wl/hermes/cmd/chat"
	"github.com/k10wl/hermes/cmd/serve"
	"github.com/k10wl/hermes/cmd/template"
	"github.com/k10wl/hermes/cmd/version"
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

func Execute(core *core.Core) error {
	rootCmd := &cobra.Command{
		Use:   "hermes",
		Short: "Tool for communication with LLM and completion instructions management",
		Long: `Host-based Extensible Response Management System
Hermes is a tool created to boost AI user experience from the terminal and browser.
It provides templating, chat persistence, and a local database.`,
		Example: ` $ hermes chat --content "Hello world!"
Hello! How can I assist you today?`,
	}

	rootCmd.AddCommand(version.CreateVersionCommand(core))
	rootCmd.AddCommand(serve.CreateServeCommand(core))
	rootCmd.AddCommand(template.CreateTemplateCommand(core))
	rootCmd.AddCommand(chat.CreateChatCommand(core))

	return rootCmd.Execute()
}
