package ai_clients

import (
	"io"
	"net/http"

	"github.com/k10wl/hermes/internal/settings"
)

type CompletionFn func(
	messages []*Message,
	parameters *Parameters,
	providers *settings.Providers,
) (*AIResponse, error)

type getter func(
	url string,
	body io.Reader,
	fillHeaders func(*http.Request) error,
) ([]byte, error)

type client interface {
	complete(
		messages []*Message,
		parameters *Parameters,
		get getter,
	) (*AIResponse, error)
}

type Message struct {
	Role    string
	Content string
}

type TokensUsage struct {
	Input  int64
	Output int64
}

type Parameters struct {
	Model       string   `json:"model"       validate:"required"`
	MaxTokens   *int64   `json:"max_tokens"`
	Temperature *float64 `json:"temperature"`
}

type AIResponse struct {
	Message
	TokensUsage TokensUsage
}

func Complete(
	messages []*Message,
	parameters *Parameters,
	providers *settings.Providers,
) (*AIResponse, error) {
	client, err := selectClient(parameters.Model, providers)
	if err != nil {
		return nil, err
	}
	return client.complete(messages, parameters, callApi)
}
