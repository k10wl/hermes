package runtime

import (
	"os"
	"path"
	"sync"
)

const (
	host = "127.0.0.1"
	port = "8123"
)

type Config struct {
	AppName   string
	Model     string
	Prompt    string
	OpenAIKey string
	ConfigDir string
	Web       bool
	Last      bool
	Host      string
	Port      string
}

var config *Config
var once sync.Once

func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		config, err = loadConfig()
	})
	return config, err
}

func loadConfig() (*Config, error) {
	var c Config
	sharedConfigDir, err := os.UserConfigDir()
	if err != nil {
		return &c, err
	}
	err = loadFlags(&c)
	if err != nil {
		return &c, err
	}
	loadEnv(&c)
	hermesConfigDir := path.Join(sharedConfigDir, c.AppName)
	err = ensureExists(hermesConfigDir)
	if err != nil {
		return &c, err
	}
	c.ConfigDir = hermesConfigDir
	return &c, nil
}
