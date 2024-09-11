package cmd

import (
	"context"

	"github.com/k10wl/hermes/cmd/chat"
	"github.com/k10wl/hermes/cmd/serve"
	"github.com/k10wl/hermes/cmd/template"
	"github.com/k10wl/hermes/cmd/version"
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hermes",
	Short: "Tool for communication with LLM and completion instructions management",
	Long: `Host-based Extensible Response Management System
hermes is a tool created to boost AI user experience from terminal and browser.
Provides templating, chat persistense and local database.`,
	Example: ` $ hermes chat --content "Hello world!"
Hello! How can I assist you today?`,
}

func init() {
	rootCmd.AddCommand(template.TemplateCommand)
	rootCmd.AddCommand(serve.ServeCommand)
	rootCmd.AddCommand(chat.ChatCommand)
	rootCmd.AddCommand(version.VersionCommand)
}

func Execute(core *core.Core) error {
	ctx := context.WithValue(context.Background(), "core", core)
	rootCmd.SetContext(ctx)
	return rootCmd.Execute()
}
