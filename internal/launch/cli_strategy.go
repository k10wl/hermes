package launch

import (
	"fmt"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
)

type cliStrategies interface {
	NewChat(*core.Core, *settings.Config) error
	LastChat(*core.Core, *settings.Config) error
	ListTemplates(*core.Core, *settings.Config) error
	UpsertTemplate(*core.Core, *settings.Config) error
	DeleteTemplate(*core.Core, *settings.Config) error
	EditTemplate(*core.Core, *settings.Config) error
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
	if countTruthyValues(
		config.ListTemplates,
		config.UpsertTemplate,
		config.DeleteTemplate,
		config.EditTemplate,
	) > 1 {
		return fmt.Errorf(
			"conflicting instruction, please review flags",
		)
	}
	if config.ListTemplates != "" {
		return l.strategies.ListTemplates(c, config)
	}
	if config.UpsertTemplate != "" {
		return l.strategies.UpsertTemplate(c, config)
	}
	if config.DeleteTemplate != "" {
		return l.strategies.DeleteTemplate(c, config)
	}
	if config.EditTemplate != "" {
		return l.strategies.EditTemplate(c, config)
	}
	if config.Last {
		return l.strategies.LastChat(c, config)
	}
	return l.strategies.NewChat(c, config)
}
