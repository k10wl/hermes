package app

import (
	"strings"

	"github.com/k10wl/hermes/internal/runtime"
)

func pickStrategy(c *runtime.Config) launchStrategy {
	if c.Web {
		return &launchWeb{}
	}
	if c.Last || c.Host != "" || c.Port != "" {
		return &launchBadInput{}
	}
	if strings.Trim(c.Prompt, " ") != "" {
		return &launchCLI{}
	}
	return &launchBadInput{}
}
