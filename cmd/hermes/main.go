package main

import (
	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/cli"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/runtime"
	"github.com/k10wl/hermes/internal/sqlite3"
	"github.com/k10wl/hermes/internal/web"
	client "github.com/k10wl/openai-client"
)

func main() {
	config, err := runtime.GetConfig()
	if err != nil {
		panic(err)
	}
	openai := ai_clients.NewOpenAIAdapter(
		client.NewOpenAIClient(config.OpenAIKey),
	)
	openai, err = openai.SetModel(config.Model)
	if err != nil {
		panic(err)
	}
	// TODO skip sqlite setup when user selectes -non-persistent
	sqlite, err := sqlite3.NewSQLite3(config)
	if err != nil {
		panic(err)
	}
	defer sqlite.Close()
	hermesCore := core.NewCore(openai, sqlite)
	cli.CLI(hermesCore, config)
	if config.Web {
		if err := web.Serve(hermesCore, config); err != nil {
			panic(err)
		}
		return
	}
}
