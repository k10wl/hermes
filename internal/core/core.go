package core

import (
	"fmt"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/db"
	"github.com/k10wl/hermes/internal/models"
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

func (c *Core) assertAI() error {
	if c == nil {
		return fmt.Errorf("core is nil")
	}
	if c.ai_client == nil {
		return fmt.Errorf("ai client is nil")
	}
	return nil
}

func messageToAIMessage(m *models.Message) ai_clients.Message {
	return ai_clients.Message{
		Content: m.Content,
		Role:    m.Role,
	}
}
