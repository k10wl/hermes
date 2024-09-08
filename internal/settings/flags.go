package settings

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func loadFlags(config *Config) {
	emptyNumericValue := -1
	template := flagStringWithShorthand(
		"template",
		"t",
		"",
		"Name of (t)emplate to be applied to provided content",
	)
	listTemplates := flagStringWithShorthand(
		"list-templates",
		"lt",
		"",
		`(L)ists (t)emplates with name that matches given regex. E.g.
`+"- `hermes -lt %`"+`- returns all templates
`+"- `hermes -lt short`"+`- return templates with name "short"`,
	)
	upsertTemplate := flagStringWithShorthand(
		"upsert-template",
		"ut",
		"",
		"Contents of (u)psert (t)emplate. Differs from golang templates, has unique delimiters. Left `--{{`, right `}}`. Is extension of golang std text/template. E.g. `--{{define \"test\"}}test--{{end}}`",
	)
	deleteTemplate := flagStringWithShorthand(
		"delete-template",
		"dt",
		"",
		`(D)eletes (t)emplate with given name. Returns error if template does not exist`,
	)
	editTemplate := flagStringWithShorthand(
		"edit-template",
		"et",
		"",
		`(E)dits (t)emplate with diven name`,
	)

	model := flagStringWithShorthand("model", "m", "gpt-4o-mini", "Completion (m)odel name")
	content := flagStringWithShorthand(
		"content",
		"c",
		"",
		`Input (c)ontent send to AI. Can contain templates (golang text/template syntax).
Interactions with other flags:
  -web  -- opens default browser in newly created chat;`,
	)
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
	temperature := flag.Float64(
		"temperature",
		1,
		"The degree of randomness or exploration in the decision-making process of an AI system.",
	)
	maxTokens := flag.Int64(
		"max-tokens",
		int64(emptyNumericValue),
		"Max tokens refer to the maximum number of words or parts of words that an AI model can use in a single output",
	)

	flag.Parse()
	config.Template = *template
	config.ListTemplates = *listTemplates
	config.UpsertTemplate = *upsertTemplate
	config.DeleteTemplate = *deleteTemplate
	config.EditTemplate = *editTemplate
	config.Content = readInput(*content)
	config.Model = *model
	config.Port = *port
	config.Web = *web
	config.Host = *host
	config.Last = *last
	config.Temperature = temperature
	if *maxTokens != int64(emptyNumericValue) {
		config.MaxTokens = maxTokens
	}
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
