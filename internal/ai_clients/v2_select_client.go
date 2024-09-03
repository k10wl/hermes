package ai_clients

import (
	"fmt"
	"strings"
	"sync"

	"github.com/k10wl/hermes/internal/settings"
)

var (
	cachedClientOpenAI *clientOpenAI
	cachedClientClaude *clientClaude
	clientMutext       sync.Mutex
)

func selectClient(input string, providers *settings.Providers) (client, error) {
	str := strings.SplitN(input, "-", 2)
	if len(str) < 2 {
		return nil, fmt.Errorf("failed to get provider %q - use gpt or claude", input)
	}
	clientMutext.Lock()
	defer clientMutext.Unlock()
	var client client
	switch provider := str[0]; provider {
	case "gpt":
		if cachedClientOpenAI == nil {
			cachedClientOpenAI = newClientOpenAI(providers.OpenAIKey)
		}
		client = cachedClientOpenAI
		break
	case "claude":
		if cachedClientClaude == nil {
			cachedClientClaude = newClientClaude(providers.AnthropicKey)
		}
		client = cachedClientClaude
		break
	default:
		return nil, fmt.Errorf("unsupported provider %q - use gpt or claude", input)
	}
	return client, nil
}
