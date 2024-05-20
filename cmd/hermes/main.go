package main

import (
	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/cli"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/runtime"
	client "github.com/k10wl/openai-client"
)

func main() {
	config := runtime.GetConfig()
	openai := ai_clients.NewOpenAIAdapter(
		client.NewOpenAIClient(*config.OpenAIKey),
	)
	openai, err := openai.SetModel(*config.Model)
	if err != nil {
		panic(err)
	}
	aiclient := core.NewCore().SetAIClient(openai)
	cli.CLI(*aiclient)
}
