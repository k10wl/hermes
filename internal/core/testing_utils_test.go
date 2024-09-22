package core_test

import (
	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/db"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
)

func createCoreAndDB() (*core.Core, db.Client) {
	db, err := sqlite3.NewSQLite3(":memory:")
	if err != nil {
		panic(err)
	}
	return core.NewCore(db, &settings.Config{}), db
}

func mockCompletion(
	messages []*ai_clients.Message,
	params *ai_clients.Parameters,
	settings *settings.Providers,
) (*ai_clients.AIResponse, error) {
	messages[0].Role = core.AssistantRole
	return &ai_clients.AIResponse{
		Message: *messages[0],
	}, nil
}
