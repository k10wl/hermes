package settings

import (
	"os"
)

func loadEnv(c *Config) {
	c.OpenAIKey = os.Getenv("HERMES_OPENAI_API_KEY")
	c.AnthropicKey = os.Getenv("HERMES_ANTHROPIC_API_KEY")
	c.DatabaseDSN = os.Getenv("HERMES_DB_DNS")
}
