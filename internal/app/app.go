package app

import (
	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/runtime"
	"github.com/k10wl/hermes/internal/sqlite3"
	client "github.com/k10wl/openai-client"
)

func Launch() error {
	config, err := runtime.GetConfig()
	if err != nil {
		return err
	}
	openai := ai_clients.NewOpenAIAdapter(
		client.NewOpenAIClient(config.OpenAIKey),
	)
	openai, err = openai.SetModel(config.Model)
	if err != nil {
		return err
	}
	sqlite, err := sqlite3.NewSQLite3(config)
	if err != nil {
		return err
	}
	defer sqlite.Close()
	hermesCore := core.NewCore(openai, sqlite)
	return pickStrategy(config).execute(hermesCore, config)
}
