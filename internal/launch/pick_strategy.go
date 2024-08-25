package launch

import (
	"strings"

	"github.com/k10wl/hermes/internal/settings"
)

func PickStrategy(c *settings.Config) launchStrategy {
	if c.Web {
		return &launchWeb{}
	}
	if c.Last || c.Host != settings.DefaultHostname || c.Port != settings.DefaultPort {
		return &launchBadInput{}
	}
	if strings.Trim(c.Content, " ") != "" {
		return &launchCLI{}
	}
	return &launchBadInput{}
}
