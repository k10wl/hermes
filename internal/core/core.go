package core

import (
	"errors"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/db"
)

const (
	UserRole      = "user"
	AssistantRole = "assistant"
	SystemRole    = "system"
)

type Core struct {
	ai_client ai_clients.AIClient
	db        db.Client
}

func NewCore(ai ai_clients.AIClient, db db.Client) *Core {
	return &Core{
		ai_client: ai,
		db:        db,
	}
}

// move to core/actions
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
