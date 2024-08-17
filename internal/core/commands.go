package core

import (
	"context"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/models"
)

type Command interface {
	Execute(context.Context) error
}

type CreateChatAndCompletionCommand struct {
	Core    *Core
	Message string
	Result  *models.Message
}

func (c *CreateChatAndCompletionCommand) Execute(ctx context.Context) error {
	if err := c.Core.assertAI(); err != nil {
		return err
	}
	chat, _, err := c.Core.db.CreateChatAndMessage(
		ctx,
		1,
		c.Message,
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
		chat.ID,
		2,
		res.Content,
	)
	c.Result = message
	return err
}

type CreateCompletionCommand struct {
	Core    *Core
	Message string
	ChatID  int64
	Result  *models.Message
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
		c.ChatID,
		1,
		c.Message,
	)
	if err != nil {
		return err
	}
	history := []ai_clients.Message{}
	for _, p := range prev {
		history = append(history, messageToAIMessage(p))
	}
	history = append(history, ai_clients.Message{Content: c.Message, Role: UserRole})
	// TODO insert used value into the db and adjust queries to receive less messages
	res, _, err := c.Core.ai_client.ChatCompletion(history)
	if err != nil {
		return err
	}
	message, err := c.Core.db.CreateMessage(
		ctx,
		c.ChatID,
		2,
		res.Content,
	)
	c.Result = message
	return err
}

type UpdateWebSettingsCommand struct {
	Core        *Core
	WebSettings models.WebSettings
}

func (c *UpdateWebSettingsCommand) Execute(ctx context.Context) error {
	return c.Core.db.UpdateWebSettings(ctx, c.WebSettings.DarkMode)
}
