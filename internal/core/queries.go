package core

import (
	"context"

	"github.com/k10wl/hermes/internal/sqlc"
)

type Query interface {
	Execute() error
}

type GetChatsQuery struct {
	Core   *Core
	Result []sqlc.Chat
}

func (c *GetChatsQuery) Execute(ctx context.Context) error {
	if err := c.Core.assertAI(); err != nil {
		return err
	}
	chats, err := c.Core.db.GetChats(ctx)
	if err != nil {
		return err
	}
	c.Result = chats
	return nil
}

type GetChatMessagesQuery struct {
	Core   *Core
	ChatID int64
	Result []sqlc.GetChatMessagesRow
}

func (c *GetChatMessagesQuery) Execute(ctx context.Context) error {
	if err := c.Core.assertAI(); err != nil {
		return err
	}
	messages, err := c.Core.db.GetChatMessages(ctx, c.ChatID)
	if err != nil {
		return err
	}
	c.Result = messages
	return nil
}
