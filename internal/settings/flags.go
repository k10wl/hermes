package settings

import (
	"flag"
	"fmt"
	"io"
	"os"

	client "github.com/k10wl/openai-client"
)

func loadFlags(config *Config) error {
	assignSharedFlags(config)
	assignTemplateFlags(config)
	assignWebFlags(config)
	flag.Parse()
	config.Content = readInput(config.Content)
	return nil
}

func assignSharedFlags(config *Config) {
	model := flagStringWithShorthand("model", "m", client.GPT3_5Turbo, "Completion (m)odel name")
	content := flagStringWithShorthand(
		"content",
		"c",
		"",
		`Input (c)ontent send to AI. Can contain templates (golang text/template syntax).
Interactions with other flags:
  -web  -- opens default browser in newly created chat;`,
	)
	config.Content = *content
	config.Model = *model
}

func assignTemplateFlags(config *Config) {
	template := flagStringWithShorthand(
		"template",
		"t",
		"",
		"Name of (t)emplate to be applied to provided content",
	)
	upsertTemplate := flagStringWithShorthand(
		"upsert-template",
		"ut",
		"",
		"Contents of (u)psert (t)emplate. Must comply golang text/template syntax",
	)
	config.Template = *template
	config.UpsertTemplate = *upsertTemplate
}

func assignWebFlags(config *Config) {
	web := flagBoolWithShorthand("web", "w", false, "Starts (w)eb server")
	last := flagBoolWithShorthand(
		"last",
		"l",
		false,
		"Opens (l)ast chat in web. Does nothing if \"-web\" was not provided",
	)
	host := flagStringWithShorthand(
		"hostname",
		"host",
		DefaultHostname,
		fmt.Sprintf("Define (host)name IP for web server. Defaults to %q", DefaultHostname),
	)
	port := flagStringWithShorthand(
		"port",
		"p",
		DefaultPort,
		fmt.Sprintf("Specify (p)ort for web server. Default to %q", DefaultPort),
	)
	config.Port = *port
	config.Web = *web
	config.Host = *host
	config.Last = *last
}

func readInput(message string) string {
	stdin, err := readStdin()
	if stdin == "" || err != nil {
		return message
	}
	return fmt.Sprintf("%s\n\n%s", stdin, message)
}

func readStdin() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}
	p, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(p), nil
}

func flagBoolWithShorthand(
	name string,
	shorthand string,
	value bool,
	usage string,
) *bool {
	var val bool
	flag.BoolVar(&val, name, value, usage)
	flag.BoolVar(&val, shorthand, value, fmt.Sprintf("shorthand for %q", name))
	return &val
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
