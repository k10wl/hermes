package test_helpers

import (
	"fmt"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
)

func MockCompletion(
	messages []*ai_clients.Message,
	params *ai_clients.Parameters,
	settings *settings.Providers,
) (*ai_clients.AIResponse, error) {
	messages[0].Role = core.AssistantRole
	messages[0].Content = fmt.Sprintf("> mocked: %s", messages[0].Content)
	return &ai_clients.AIResponse{
		Message: *messages[0],
	}, nil
}
