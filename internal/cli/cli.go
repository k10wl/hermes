package cli

import (
	"context"
	"fmt"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
)

func NewChat(c *core.Core, config *settings.Config) error {
	sendMessage := core.NewCreateChatAndCompletionCommand(
		c,
		core.UserRole,
		config.Content,
		config.Template,
	)
	if err := sendMessage.Execute(context.Background()); err != nil {
		return err
	}
	fmt.Fprintf(config.Stdoout, "%s\n", sendMessage.Result.Content)
	return nil
}

func LastChat(c *core.Core, config *settings.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queryChats := core.LatestChatQuery{Core: c}
	if err := queryChats.Execute(ctx); err != nil {
		return err
	}
	completionCommand := core.NewCreateCompletionCommand(
		c,
		queryChats.Result.ID,
		core.UserRole,
		config.Content,
		config.Template,
	)
	if err := completionCommand.Execute(ctx); err != nil {
		return err
	}
	fmt.Fprintf(config.Stdoout, "%s\n", completionCommand.Result.Content)
	return nil
}

func UpsertTemplate(c *core.Core, config *settings.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return core.NewUpsertTemplateCommand(c, config.UpsertTemplate).Execute(ctx)
}
