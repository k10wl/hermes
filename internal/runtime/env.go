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

func ensureExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
