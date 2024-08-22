package app

import (
	"strings"

	"github.com/k10wl/hermes/internal/settings"
)

func PickStrategy(c *settings.Config) launchStrategy {
	if c.Web {
		return &launchWeb{}
	}
	if c.Last || c.Host != settings.DefaultHost || c.Port != settings.DefaultPort {
		return &launchBadInput{}
	}
	if strings.Trim(c.Prompt, " ") != "" {
		return &launchCLI{}
	}
	return &launchBadInput{}
}
