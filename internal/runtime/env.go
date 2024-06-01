package runtime

import "os"

var appName = "hermes" // fallback name
func loadEnv(c *Config) {
	name, ok := os.LookupEnv("APP_NAME")
	if ok {
		c.AppName = name
	} else {
		c.AppName = appName
	}
	openAIKey := os.Getenv("OPEN_AI_KEY")
	c.OpenAIKey = openAIKey
}
