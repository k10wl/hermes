package settings

import (
	"context"
	"io"
	"os"
	"path"
	"sync"
)

const Version = "4.0.0"
const DefaultHostname = "127.0.0.1"

var DefaultPort = "8123"            // changes in ldflag for dev mode
var DefaultDatabaseName = "main.db" // changes in ldflag for dev mode
var appName = "hermes"              // changes in ldflag for dev mode
var config *Config
var once sync.Once

type Config struct {
	Settings
	Providers
	TemplateFlags
	CLIFlags
	WebFlags
}

type Settings struct {
	Version         string
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

type CLIFlags struct {
	Model       string
	Content     string
	Last        bool
	MaxTokens   *int64
	Temperature *float64
}

type WebFlags struct {
	Web  bool
	Host string
	Port string
}

type TemplateFlags struct {
	Template       string
	ListTemplates  string
	UpsertTemplate string
	DeleteTemplate string
	EditTemplate   string
}

func GetConfig(stdin io.Reader, stdout io.Writer, stderr io.Writer) (*Config, error) {
	var err error
	once.Do(func() {
		config, err = loadConfig(stdin, stdout, stderr)
	})
	return config, err
}

func loadConfig(stdin io.Reader, stdout io.Writer, stderr io.Writer) (*Config, error) {
	var c Config
	c.AppName = appName
	loadEnv(&c)
	if err := prepareDNS(&c); err != nil {
		return &c, err
	}
	c.Version = Version
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
	c.DatabaseDSN = path.Join(hermesConfigDir, "main.db")
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
