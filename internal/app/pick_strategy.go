package app

import (
	"strings"

	"github.com/k10wl/hermes/internal/runtime"
)

func PickStrategy(c *runtime.Config) launchStrategy {
	if c.Web {
		return &launchWeb{}
	}
	if c.Last || c.Host != runtime.DefaultHost || c.Port != runtime.DefaultPort {
		return &launchBadInput{}
	}
	if strings.Trim(c.Prompt, " ") != "" {
		return &launchCLI{}
	}
	return &launchBadInput{}
}
