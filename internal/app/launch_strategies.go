package app

import (
	"context"
	"fmt"

	"github.com/k10wl/hermes/internal/cli"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/runtime"
	"github.com/k10wl/hermes/internal/web"
)

type launchStrategy interface {
	execute(*core.Core, *runtime.Config) error
}

type launchWeb struct{}

func (l *launchWeb) execute(c *core.Core, config *runtime.Config) error {
	if config.Prompt != "" {
		sendMessage := core.CreateChatAndCompletionCommand{
			Core:    c,
			Message: config.Prompt,
		}
		if err := sendMessage.Execute(context.Background()); err != nil {
			return err
		}
	}
	return web.Serve(c, config)
}

type launchCLI struct{}

func (l *launchCLI) execute(c *core.Core, config *runtime.Config) error {
	return cli.CLI(c, config)
}

type launchBadInput struct{}

func (l *launchBadInput) execute(c *core.Core, config *runtime.Config) error {
	fmt.Println(cli.HelpString)
	return nil
}
