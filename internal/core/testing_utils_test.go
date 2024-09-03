package core_test

import (
	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/db"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/sqlite3"
)

type MockAIClient struct{}

func (mockClient MockAIClient) ChatCompletion(
	messages []ai_clients.Message,
) (ai_clients.Message, int, error) {
	messages[0].Role = core.AssistantRole
	return messages[0], 1, nil
}

func createCoreAndDB() (*core.Core, db.Client) {
	db, err := sqlite3.NewSQLite3(
		&settings.Config{Settings: settings.Settings{DatabaseDSN: ":memory:"}},
	)
	if err != nil {
		panic(err)
	}
	return core.NewCore(MockAIClient{}, db), db
}
