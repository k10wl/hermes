package core

import (
	"errors"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
)

const (
	UserRole      = "user"
	AssistantRole = "assistant"
	SystemRole    = "system"
)

type Core struct {
	ai_client ai_clients.AIClient
}

// core should know how to connect database with ai clients
// I don't want to have API keys burned down into binary, they should be passed
// as ENV variables. The application will only take the name of the ENV that
// holds value for the actual API keys

func NewCore() *Core {
	return &Core{}
}

func (c *Core) SetAIClient(a ai_clients.AIClient) *Core {
	c.ai_client = a
	return c
}

func (c *Core) SendMessage(message string) (string, error) {
	if c.ai_client == nil {
		return "", errors.New("ai client not set")
	}
	res, err := c.ai_client.ChatCompletion([]ai_clients.Message{{Content: message, Role: UserRole}})
	if err != nil {
		return "", err
	}
	return res.Content, nil
}
