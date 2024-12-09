package ai_clients

import (
	"fmt"
	"time"
)

type mock struct{}

func (m mock) complete(
	messages []*Message,
	parameters *Parameters,
	get getter,
) (*AIResponse, error) {
	messages[len(messages)-1].Role = "assistant"
	messages[len(messages)-1].Content = fmt.Sprintf(
		"> mocked: %s",
		messages[len(messages)-1].Content,
	)
	time.Sleep(time.Second)
	return &AIResponse{
		Message: *messages[len(messages)-1],
	}, nil
}
