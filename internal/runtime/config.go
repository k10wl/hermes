package runtime

import (
	"context"
	"io"
	"os"
	"path"
	"sync"
)

const (
	DefaultHost = "127.0.0.1"
	DefaultPort = "8123"
)

type Config struct {
	AppName         string
	Model           string
	Prompt          string
	OpenAIKey       string
	ConfigDir       string
	Web             bool
	Last            bool
	Host            string
	Port            string
	Stdin           io.Reader
	Stdoout         io.Writer
	Stderr          io.Writer
	DatabaseDSN     string
	ShutdownContext context.Context
}

var config *Config
var once sync.Once

func GetConfig(stdin io.Reader, stdout io.Writer, stderr io.Writer) (*Config, error) {
	var err error
	once.Do(func() {
		config, err = loadConfig(stdin, stdout, stderr)
	})
	return config, err
}

func loadConfig(stdin io.Reader, stdout io.Writer, stderr io.Writer) (*Config, error) {
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
	c.ConfigDir = hermesConfigDir
	c.DatabaseDSN = path.Join(hermesConfigDir, "main.db")
	if c.DatabaseDSN != ":memory:" {
		err = ensureExists(hermesConfigDir)
		if err != nil {
			return &c, err
		}
	}
	c.Stdin = stdin
	c.Stdoout = stdout
	c.Stderr = stderr
	c.ShutdownContext = context.Background()
	return &c, nil
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
