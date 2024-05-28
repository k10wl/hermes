package core

import (
	"context"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/sqlc"
)

type Command interface {
	Execute(context.Context) error
}

type CreateChatAndCompletionCommand struct {
	Core    *Core
	Message string
	Result  sqlc.Message
}

func (c *CreateChatAndCompletionCommand) Execute(ctx context.Context) error {
	if err := c.Core.assertAI(); err != nil {
		return err
	}
	chat, _, err := c.Core.db.CreateChatAndMessage(
		ctx,
		sqlc.CreateMessageParams{Content: c.Message, RoleID: 1},
	)
	if err != nil {
		return err
	}
	// TODO insert used value into the db and adjust queries to receive less messages
	res, _, err := c.Core.ai_client.ChatCompletion(
		[]ai_clients.Message{{Content: c.Message, Role: UserRole}},
	)
	if err != nil {
		return err
	}
	message, err := c.Core.db.CreateMessage(
		ctx,
		sqlc.CreateMessageParams{ChatID: chat.ID, Content: res.Content, RoleID: 2},
	)
	if err != nil {
		return err
	}
	c.Result = message
	return nil
}

type CreateCompletionCommand struct {
	Core    *Core
	Message string
	ChatID  int64
	Result  sqlc.Message
}

func (c *CreateCompletionCommand) Execute(ctx context.Context) error {
	if err := c.Core.assertAI(); err != nil {
		return err
	}
	prev, err := c.Core.db.GetChatMessages(ctx, c.ChatID)
	if err != nil {
		return err
	}
	_, err = c.Core.db.CreateMessage(
		ctx,
		sqlc.CreateMessageParams{ChatID: c.ChatID, Content: c.Message, RoleID: 1},
	)
	if err != nil {
		return err
	}
	history := []ai_clients.Message{}
	for _, p := range prev {
		history = append(history, sqlcMessageToAIMessage(p))
	}
	history = append(history, ai_clients.Message{Content: c.Message, Role: UserRole})
	// TODO insert used value into the db and adjust queries to receive less messages
	res, _, err := c.Core.ai_client.ChatCompletion(history)
	if err != nil {
		return err
	}
	message, err := c.Core.db.CreateMessage(
		ctx,
		sqlc.CreateMessageParams{ChatID: c.ChatID, Content: res.Content, RoleID: 2},
	)
	if err != nil {
		return err
	}
	c.Result = message
	return nil
}
