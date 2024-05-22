package core

import (
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
