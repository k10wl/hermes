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
	prompt := flag.String("prompt", "", "prompt for instant completion")
	sufix := flag.String("sufix", "", "sufix for prompt (useful with vim)")
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
	*c.Prompt = fmt.Sprintf("%s\n\n%s", *c.Prompt, *sufix)
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
