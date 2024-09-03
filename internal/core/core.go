package core

import (
	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/db"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/settings"
)

const (
	UserRole      = "user"
	AssistantRole = "assistant"
	SystemRole    = "system"
)

type Core struct {
	db     db.Client
	config *settings.Config
}

func NewCore(db db.Client, config *settings.Config) *Core {
	return &Core{
		db:     db,
		config: config,
	}
}

func messageToAIMessage(m *models.Message) *ai_clients.Message {
	return &ai_clients.Message{
		Content: m.Content,
		Role:    m.Role,
	}
}
