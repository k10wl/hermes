package ai_clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k10wl/hermes/internal/ai_clients/openai"
)

type clientOpenAI struct {
	apiKey string
	apiUrl string
}

func newClientOpenAI(apiKey string) *clientOpenAI {
	return &clientOpenAI{
		apiKey: apiKey,
		apiUrl: "https://api.openai.com/v1/chat/completions",
	}
}

func (client clientOpenAI) complete(
	messages []*Message,
	parameters *Parameters,
	get getter,
) (*AIResponse, error) {
	data, err := client.prepare(messages, parameters)
	if err != nil {
		return nil, err
	}
	res, err := get(
		client.apiUrl,
		bytes.NewBuffer(data),
		client.fillHeaders,
	)
	if err != nil {
		return nil, err
	}
	var openaiResponse openai.ChatCompletionResponse
	err = json.Unmarshal(res, &openaiResponse)
	if err != nil {
		return nil, err
	}
	return client.decodeResponse(&openaiResponse)
}

func (client clientOpenAI) prepare(messages []*Message, parameters *Parameters) ([]byte, error) {
	data := openai.ChatCompletionRequest{
		Model:    parameters.Model,
		Messages: client.encodeMessages(messages),
	}
	if parameters.MaxTokens != nil {
		data.MaxTokens = *parameters.MaxTokens
	}
	if parameters.Temperature != nil {
		data.Temperature = *parameters.Temperature
	}
	return json.Marshal(data)
}

func (client clientOpenAI) encodeMessages(messages []*Message) []*openai.Message {
	result := make([]*openai.Message, len(messages))
	for i, v := range messages {
		result[i] = &openai.Message{
			Role:    v.Role,
			Content: v.Content,
		}
	}
	return result
}

func (client clientOpenAI) decodeMessage(messages openai.Message) *Message {
	return &Message{
		Content: messages.Content,
		Role:    messages.Role,
	}
}

func (client clientOpenAI) decodeResponse(
	response *openai.ChatCompletionResponse,
) (*AIResponse, error) {
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	messages := []*Message{}
	for _, val := range response.Choices {
		messages = append(messages, client.decodeMessage(val.Message))
	}
	return &AIResponse{
		Message: Message{
			Content: response.Choices[0].Message.Content,
			Role:    response.Choices[0].Message.Role,
		},
		TokensUsage: TokensUsage{
			Input:  response.Usage.PromptTokens,
			Output: response.Usage.CompletionTokens,
		},
	}, nil
}

func (client clientOpenAI) fillHeaders(req *http.Request) error {
	if client.apiKey == "" {
		return fmt.Errorf("OpenAI API key was not provided\n")
	}
	req.Header.Add("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return nil
}
