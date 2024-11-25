package settings

import (
	"os"
)

const (
	HermesOpenAIApiKeyName    = "HERMES_OPENAI_API_KEY"
	HermesAnthropicApiKeyName = "HERMES_ANTHROPIC_API_KEY"
)

func loadEnv(c *Config) {
	c.OpenAIKey = os.Getenv(HermesOpenAIApiKeyName)
	c.AnthropicKey = os.Getenv(HermesAnthropicApiKeyName)
	c.DatabaseDSN = os.Getenv("HERMES_DB_DNS")
}
