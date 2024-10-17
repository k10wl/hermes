package core_test

import (
	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
)

func mockCompletion(
	messages []*ai_clients.Message,
	params *ai_clients.Parameters,
	settings *settings.Providers,
) (*ai_clients.AIResponse, error) {
	messages[0].Role = core.AssistantRole
	return &ai_clients.AIResponse{
		Message: *messages[0],
	}, nil
}
