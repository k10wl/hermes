package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
)

type CLIStrategies struct{}

func (cli *CLIStrategies) NewChat(c *core.Core, config *settings.Config) error {
	sendMessage := core.NewCreateChatAndCompletionCommand(
		c,
		core.UserRole,
		config.Content,
		config.Template,
		&ai_clients.Parameters{
			Model: config.Model,
		},
		ai_clients.Complete,
	)
	if err := sendMessage.Execute(context.Background()); err != nil {
		return err
	}
	fmt.Fprintf(config.Stdoout, "%s\n", sendMessage.Result.Content)
	return nil
}

func (cli *CLIStrategies) LastChat(c *core.Core, config *settings.Config) error {
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
		&ai_clients.Parameters{
			Model: config.Model,
		},
		ai_clients.Complete,
	)
	if err := completionCommand.Execute(ctx); err != nil {
		return err
	}
	fmt.Fprintf(config.Stdoout, "%s\n", completionCommand.Result.Content)
	return nil
}

func (cli *CLIStrategies) UpsertTemplate(c *core.Core, config *settings.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return core.NewUpsertTemplateCommand(c, config.UpsertTemplate).Execute(ctx)
}

func (cli *CLIStrategies) ListTemplates(c *core.Core, config *settings.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	getTemplatesByRegexp := core.NewGetTemplatesByRegexp(c, config.ListTemplates)
	if err := getTemplatesByRegexp.Execute(ctx); err != nil {
		return err
	}
	return listTemplates(getTemplatesByRegexp.Result, config.Stdoout)
}

func (cli *CLIStrategies) DeleteTemplate(c *core.Core, config *settings.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := core.NewDeleteTemplateByName(c, config.DeleteTemplate).Execute(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(config.Stdoout, "Successfully deleted %q\n", config.DeleteTemplate)
	return nil
}

func (cli *CLIStrategies) EditTemplate(c *core.Core, config *settings.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	query := core.NewGetTemplatesByNamesQuery(c, []string{config.EditTemplate})
	if err := query.Execute(ctx); err != nil {
		return err
	}
	if len(query.Result) != 1 {
		return fmt.Errorf("failed to get template\n")
	}
	res, err := OpenInEditor(query.Result[0].Content, os.Stdin, config.Stdoout, config.Stderr)
	if err != nil {
		return err
	}
	if res == query.Result[0].Content {
		fmt.Fprintf(config.Stdoout, "contents are identical\n")
		return nil
	}
	err = core.NewEditTemplateByName(c, config.EditTemplate, res).Execute(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(config.Stdoout, "successfully edited template\n")
	return nil
}
