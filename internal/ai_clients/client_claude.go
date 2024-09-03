package ai_clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/k10wl/hermes/internal/ai_clients/claude"
)

type clientClaude struct {
	apiKey            string
	apiUrl            string
	maxTokensFallback int64
}

func newClientClaude(apiKey string) *clientClaude {
	return &clientClaude{
		apiKey:            apiKey,
		apiUrl:            "https://api.anthropic.com/v1/messages",
		maxTokensFallback: 4096,
	}
}

var claudeModelsShorthands = map[string]string{
	"claude-sonet": "claude-3-5-sonnet-20240620",
	"claude-opus":  "claude-3-opus-20240229",
	"claude-haiku": "claude-3-haiku-20240307",
}

func (client clientClaude) complete(
	messages []*Message,
	parameters *Parameters,
	get getter,
) (*AIResponse, error) {
	data, err := client.prepare(messages, parameters)
	if err != nil {
		return nil, err
	}
	data, err = get(
		client.apiUrl,
		bytes.NewReader(data),
		client.fillHeaders,
	)
	if err != nil {
		return nil, err
	}
	var response claude.MessagesResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return client.decodeResult(&response)
}

func (client clientClaude) prepare(
	messages []*Message,
	parameters *Parameters,
) ([]byte, error) {
	model := parameters.Model
	if fullName, ok := claudeModelsShorthands[model]; ok {
		model = fullName
	}
	encodedMessages, systemPrompt := client.encodeMessages(messages)
	data := claude.MessagesRequest{
		Model:     model,
		Messages:  encodedMessages,
		System:    systemPrompt,
		MaxTokens: client.maxTokensFallback,
	}
	if parameters.MaxTokens != nil {
		data.MaxTokens = *parameters.MaxTokens
	}
	if parameters.Temperature != nil {
		data.Temperature = *parameters.Temperature
	}
	return json.Marshal(data)
}

func (client clientClaude) encodeMessages(
	messages []*Message,
) ([]*claude.MessageContent, string) {
	result := make([]*claude.MessageContent, len(messages))
	sb := &strings.Builder{}
	for i, v := range messages {
		if v.Role == "system" {
			sb.WriteString(v.Content + "\n")
			continue
		}
		result[i] = &claude.MessageContent{
			Role:    v.Role,
			Content: v.Content,
		}
	}
	return result, sb.String()
}

func (client clientClaude) decodeResult(response *claude.MessagesResponse) (*AIResponse, error) {
	if len(response.Content) == 0 {
		return nil, fmt.Errorf("empty response content")
	}
	return &AIResponse{
		Message: Message{
			Content: response.Content[0].Text,
			Role:    response.Role,
		},
		TokensUsage: TokensUsage{
			Input:  response.Usage.InputTokens,
			Output: response.Usage.OutputTokens,
		},
	}, nil
}

func (client clientClaude) fillHeaders(r *http.Request) error {
	if client.apiKey == "" {
		return fmt.Errorf("Claude API key was not provided\n")
	}
	r.Header.Set("x-api-key", client.apiKey)
	r.Header.Set("anthropic-version", "2023-06-01")
	r.Header.Set("content-type", "application/json")
	return nil
}
