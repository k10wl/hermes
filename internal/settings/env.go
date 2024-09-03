package settings

import (
	"os"
)

func loadEnv(c *Config) {
	openAIKey := os.Getenv("HERMES_OPENAI_API_KEY")
	anthropicKey := os.Getenv("HERMES_ANTHROPIC_API_KEY")
	c.OpenAIKey = openAIKey
	c.AnthropicKey = anthropicKey
}
