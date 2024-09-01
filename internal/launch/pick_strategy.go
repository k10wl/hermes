package launch

import (
	"github.com/k10wl/hermes/internal/cli"
	"github.com/k10wl/hermes/internal/settings"
)

func PickStrategy(config *settings.Config) launchStrategy {
	if config.Web {
		return &launchWeb{}
	}
	if countTruthyValues(
		config.Content,
		config.UpsertTemplate,
		config.ListTemplates,
		config.Template,
		config.DeleteTemplate,
		config.EditTemplate,
	) != 0 {
		return newLaunchCLI(&cli.CLIStrategies{})
	}
	return &launchBadInput{}
}
