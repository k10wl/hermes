package ai_clients

import (
	"errors"

	client "github.com/k10wl/openai-client"
)

type OpenAIClient interface {
}

type OpenAICompletionModel interface {
}

type OpenAIAdapter struct {
	model  *client.ChatCompletionModel
	client *client.OpenAIClient
}

func NewOpenAIAdapter(client *client.OpenAIClient) *OpenAIAdapter {
	return &OpenAIAdapter{client: client}
}

func (a *OpenAIAdapter) ChatCompletion(message []Message) (Message, error) {
	var res Message
	if a.model == nil {
		return res, errors.New("model was not provided")
	}
	history := []client.Message{}
	for _, v := range message {
		history = append(history, a.messageEncoder(v))
	}
	c, err := a.model.ChatCompletion(history)
	if err != nil {
		return res, err
	}
	return a.messageDecoder(c.Choices[0].Message), nil
}

func (a *OpenAIAdapter) SetModel(model string) (*OpenAIAdapter, error) {
	m, err := client.NewChatCompletionModel(a.client, model)
	if err != nil {
		return nil, err
	}
	a.model = m
	return a, nil
}

func (a *OpenAIAdapter) messageDecoder(message client.Message) Message {
	return Message{
		Role:    message.Role,
		Content: message.Content,
	}
}

func (a *OpenAIAdapter) messageEncoder(message Message) client.Message {
	return client.Message{
		Role:    message.Role,
		Content: message.Content,
	}
}
