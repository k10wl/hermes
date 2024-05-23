package core

import (
	"context"
	"fmt"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/sqlc"
)

type Command interface {
	Execute() error
}

type CreateChatAndCompletionCommand struct {
	Core    *Core
	Message string
	Result  string
}

func (c *CreateChatAndCompletionCommand) Execute(ctx context.Context) error {
	if c.Core == nil || c.Core.ai_client == nil {
		return fmt.Errorf("ai client not set")
	}
	chat, _, err := c.Core.db.CreateChatAndMessage(
		ctx,
		sqlc.CreateMessageParams{Content: c.Message, RoleID: 1},
	)
	if err != nil {
		return err
	}
	res, err := c.Core.ai_client.ChatCompletion(
		[]ai_clients.Message{{Content: c.Message, Role: UserRole}},
	)
	if err != nil {
		return err
	}
	_, err = c.Core.db.CreateMessage(
		ctx,
		sqlc.CreateMessageParams{ChatID: chat.ID, Content: res.Content, RoleID: 2},
	)
	if err != nil {
		return err
	}
	c.Result = res.Content
	return nil
}
