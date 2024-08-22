package app

import (
	"context"
	"fmt"

	"github.com/k10wl/hermes/internal/cli"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/web"
)

type launchStrategy interface {
	Execute(*core.Core, *settings.Config) error
}

type launchWeb struct{}

func (l *launchWeb) Execute(c *core.Core, config *settings.Config) error {
	if config.Input != "" {
		sendMessage := core.CreateChatAndCompletionCommand{
			Core:    c,
			Message: config.Input,
			Role:    "user",
		}
		if err := sendMessage.Execute(context.Background()); err != nil {
			return err
		}
	}
	web.OpenBrowser(web.GetUrl(fmt.Sprintf("%s:%s", config.Host, config.Port), c, config))
	return web.Serve(c, config)
}

type launchCLI struct{}

func (l *launchCLI) Execute(c *core.Core, config *settings.Config) error {
	return cli.CLI(c, config)
}

type launchBadInput struct{}

func (l *launchBadInput) Execute(c *core.Core, config *settings.Config) error {
	fmt.Fprintf(config.Stdoout, "%s\n", cli.HelpString)
	return nil
}
