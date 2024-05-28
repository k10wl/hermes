package runtime

import (
	"flag"
	"fmt"
	"io"
	"os"

	client "github.com/k10wl/openai-client"
)

func loadFlags(c *Config) error {
	// FIXME this should be loaded from some kind of user preferences
	model := flag.String("model", client.GPT3_5Turbo, "ai model name")
	message := flagStringWithShorthand(
		"message",
		"m",
		"",
		"Inline prompt message attached to end of Stdin string, or used as standalone prompt string",
	)
	web := flag.Bool("web", false, "Starts web server")
	host := flagStringWithShorthand(
		"host",
		"h",
		host,
		"Host for web server. Optional, does nothing if \"-web\" was not set",
	)
	port := flagStringWithShorthand(
		"port",
		"p",
		port,
		"Port for web server. Optional, does nothing if \"-web\" was not set",
	)
	flag.Parse()
	c.Model = *model
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		p, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		c.Prompt = string(p)
	}
	if *message != "" {
		if c.Prompt == "" {
			c.Prompt = *message
		} else {
			c.Prompt = fmt.Sprintf("%s\n\n%s", c.Prompt, *message)
		}
	}
	c.Host = *host
	c.Port = *port
	c.Web = *web
	return nil
}

func flagStringWithShorthand(
	name string,
	shorthand string,
	value string,
	usage string,
) *string {
	var val string
	flag.StringVar(&val, name, value, usage)
	flag.StringVar(&val, shorthand, value, fmt.Sprintf("shorthand for %q", name))
	return &val
}
