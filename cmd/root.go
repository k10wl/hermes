package cmd

import (
	"context"

	"github.com/k10wl/hermes/cmd/chat"
	"github.com/k10wl/hermes/cmd/serve"
	"github.com/k10wl/hermes/cmd/template"
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(0),
	Use:   "hermes",
	Short: "tool for productive communication with and management of AI",
	Long: `Host-based Extensible Response Management System
hermes is a tool created to boost AI user experience from terminal and browser.
Provies templating, chat persistense, local database and more.
Offers access to OpenAI and Anthropic (Claude) completion API by personal keys.

Example:
    $ HERMES_OPENAI_API_KEY=your-own-key hermes chat <<< "who are you?"
    I am a language model AI designed to assist with answering questions and
    providing information to the best of my knowledge and abilities.`,
}

func init() {
	rootCmd.AddCommand(template.TemplateCommand)
	rootCmd.AddCommand(serve.ServeCommand)
	rootCmd.AddCommand(chat.ChatCommand)
}

func Execute(core *core.Core) error {
	ctx := context.WithValue(context.Background(), "core", core)
	rootCmd.SetContext(ctx)
	return rootCmd.Execute()
}
