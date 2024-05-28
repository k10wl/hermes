package ai_clients

import (
	"errors"
	"slices"

	client "github.com/k10wl/openai-client"
	"github.com/tiktoken-go/tokenizer"
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

func (a *OpenAIAdapter) ChatCompletion(messages []Message) (Message, int, error) {
	var res Message
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return res, 0, err
	}
	if a.model == nil {
		return res, 0, errors.New("model was not provided")
	}
	history := []client.Message{}
	usedMessages := 0
	tokens := 0
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		_, t, _ := enc.Encode(a.messageEncoder(message).Content)
		tokens += len(t)
		if tokens > a.model.TokenLimit {
			break
		}
		history = slices.Insert(history, 0, a.messageEncoder(message))
		usedMessages++
	}
	c, err := a.model.ChatCompletion(history)
	if err != nil {
		return res, 0, err
	}
	return a.messageDecoder(c.Choices[0].Message), usedMessages, nil
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
