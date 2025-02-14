package ai_clients

import (
	"fmt"
	"strings"
	"sync"

	"github.com/k10wl/hermes/internal/settings"
)

var (
	cachedClientOpenAI    *clientOpenAI
	cachedClientAnthropic *clientClaude
	clientMutext          sync.Mutex
)

func selectClient(provider string, providers *settings.Providers) (client, error) {
	if config, err := settings.GetInstance(); err == nil && config.MockCompletion {
		return mock{}, nil
	}
	clientMutext.Lock()
	defer clientMutext.Unlock()
	var client client
	switch provider {
	case "openai":
		if cachedClientOpenAI == nil {
			cachedClientOpenAI = newClientOpenAI(providers.OpenAIKey)
		}
		client = cachedClientOpenAI
	case "anthropic":
		if cachedClientAnthropic == nil {
			cachedClientAnthropic = newClientClaude(providers.AnthropicKey)
		}
		client = cachedClientAnthropic
	default:
		return nil, fmt.Errorf("unsupported provider %q - use openai/model or anthropic/model", provider)
	}
	return client, nil
}

func extractProviderAndModel(input string) (string, string, error) {
	str := strings.SplitN(input, "/", 2)
	if len(str) < 2 {
		return "", "", fmt.Errorf("failed to get provider and model from %q - expected format: provider/model", input)
	}
	return str[0], str[1], nil
}
