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
	// TEST
	web := flag.Bool("web", false, "Starts web server")
	last := flag.Bool(
		"last",
		false,
		"Opens last chat in web. Optional, does nothing if \"-web\" was not provided",
	)
	host := flagStringWithShorthand(
		"host",
		"h",
		DefaultHost,
		"Host for web server. Optional, does nothing if \"-web\" was not provided",
	)
	port := flagStringWithShorthand(
		"port",
		"p",
		DefaultPort,
		"Port for web server. Optional, does nothing if \"-web\" was not provided",
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
	c.Last = *last
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
