package settings

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
)

const Version = "4.6.1"
const VersionDate = "2024-12-02"

const DefaultDatabaseName = "main.db"

var appName = "hermes" // changes in ldflag for dev mode
var config *Config
var once sync.Once

type Config struct {
	Settings
	Providers
}

type Settings struct {
	Version         string
	VersionDate     string
	AppName         string
	ConfigDir       string
	DatabaseDSN     string
	ShutdownContext context.Context
	Stdin           io.Reader
	Stdoout         io.Writer
	Stderr          io.Writer
}

type Providers struct {
	OpenAIKey    string
	AnthropicKey string
}

func GetConfig(stdin io.Reader, stdout io.Writer, stderr io.Writer) (*Config, error) {
	var err error
	once.Do(func() {
		config, err = loadConfig(stdin, stdout, stderr)
	})
	return config, err
}

func GetInstance() (*Config, error) {
	if config == nil {
		return nil, fmt.Errorf("config was not initialized")
	}
	return config, nil
}

func loadConfig(stdin io.Reader, stdout io.Writer, stderr io.Writer) (*Config, error) {
	var c Config
	c.AppName = appName
	loadEnv(&c)
	if err := prepareDNS(&c); err != nil {
		return &c, err
	}
	c.Version = Version
	c.VersionDate = VersionDate
	c.Stdin = stdin
	c.Stdoout = stdout
	c.Stderr = stderr
	c.ShutdownContext = context.Background()
	return &c, nil
}

func prepareDNS(c *Config) error {
	if c.DatabaseDSN == ":memory:" || c.DatabaseDSN != "" {
		return nil
	}
	sharedConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	hermesConfigDir := path.Join(sharedConfigDir, c.AppName)
	c.ConfigDir = hermesConfigDir
	c.DatabaseDSN = path.Join(hermesConfigDir, DefaultDatabaseName)
	err = ensureExists(hermesConfigDir)
	return err
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
