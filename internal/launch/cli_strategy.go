package launch

import (
	"fmt"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
)

type cliStrategies interface {
	NewChat(*core.Core, *settings.Config) error
	LastChat(*core.Core, *settings.Config) error
	ViewTemplates(*core.Core, *settings.Config) error
	UpsertTemplate(*core.Core, *settings.Config) error
}

type launchCLI struct {
	strategies cliStrategies
}

func newLaunchCLI(strategies cliStrategies) *launchCLI {
	return &launchCLI{
		strategies: strategies,
	}
}

func (l *launchCLI) Execute(c *core.Core, config *settings.Config) error {
	if config.ViewTemplates != "" && config.UpsertTemplate != "" {
		return fmt.Errorf(
			"conflicting instruction, please do not combine view templates with upsert template",
		)
	}
	if config.ViewTemplates != "" {
		return l.strategies.ViewTemplates(c, config)
	}
	if config.UpsertTemplate != "" {
		return l.strategies.UpsertTemplate(c, config)
	}
	if config.Last {
		return l.strategies.LastChat(c, config)
	}
	return l.strategies.NewChat(c, config)
}
