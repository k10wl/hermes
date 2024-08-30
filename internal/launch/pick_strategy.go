package launch

import (
	"strings"

	"github.com/k10wl/hermes/internal/cli"
	"github.com/k10wl/hermes/internal/settings"
)

func PickStrategy(config *settings.Config) launchStrategy {
	if config.Web {
		return &launchWeb{}
	}
	if strings.Trim(config.Content, " ") != "" ||
		config.UpsertTemplate != "" ||
		config.Template != "" {
		return newLaunchCLI(&cli.CLIStrategies{})
	}
	if config.Last || config.Host != settings.DefaultHostname ||
		config.Port != settings.DefaultPort {
		return &launchBadInput{}
	}
	return &launchBadInput{}
}
