package runtime

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	client "github.com/k10wl/openai-client"
)

type Config struct {
	AppName   *string
	Model     *string
	Prompt    *string
	OpenAIKey *string
}

var config *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		config = loadConfig()
	})
	return config
}

func loadConfig() *Config {
	var c Config
	loadFlags(&c)
	loadEnv(&c)
	return &c
}

func loadFlags(c *Config) {
	model := flag.String("model", client.GPT3_5Turbo, "ai model name")
	prompt := flagStringWithShorthand(
		"prompt",
		"p",
		"",
		"prompt for instant completion",
	)
	message := flagStringWithShorthand(
		"message",
		"m",
		"",
		"message attached to the end of the prompt (useful with vim)",
	)
	flag.Parse()
	c.Model = model
	c.Prompt = prompt
	if *c.Prompt == "" {
		p, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		*c.Prompt = string(p)
	}
	if *message != "" {
		*c.Prompt = fmt.Sprintf("%s\n\n%s", *c.Prompt, *message)
	}
}

func flagStringWithShorthand(
	name string,
	shorthand string,
	value string,
	usage string,
) *string {
	var val string
	flag.StringVar(&val, name, value, usage)
	if val != "" {
		return &val
	}
	flag.StringVar(&val, shorthand, value, usage)
	return &val
}

var appName = "hermes" // fallback name
func loadEnv(c *Config) {
	name, ok := os.LookupEnv("APP_NAME")
	if ok {
		c.AppName = &name
	} else {
		c.AppName = &appName
	}
	openAIKey := os.Getenv("OPEN_AI_KEY")
	c.OpenAIKey = &openAIKey
}
