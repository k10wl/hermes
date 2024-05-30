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

func (q *GetChatsQuery) Execute(ctx context.Context) error {
	if err := q.Core.assertAI(); err != nil {
		return err
	}
	chats, err := q.Core.db.GetChats(ctx)
	if err != nil {
		return err
	}
	q.Result = chats
	return nil
}

type GetChatMessagesQuery struct {
	Core   *Core
	ChatID int64
	Result []sqlc.GetChatMessagesRow
}

func (q *GetChatMessagesQuery) Execute(ctx context.Context) error {
	if err := q.Core.assertAI(); err != nil {
		return err
	}
	messages, err := q.Core.db.GetChatMessages(ctx, q.ChatID)
	if err != nil {
		return err
	}
	q.Result = messages
	return nil
}

type WebSettingsQuery struct {
	Core   *Core
	Result sqlc.WebSetting
}

func (q *WebSettingsQuery) Execute(ctx context.Context) error {
	if err := q.Core.assertAI(); err != nil {
		return err
	}
	setting, err := q.Core.db.GetWebSettings(ctx)
	q.Result = setting
	return err
}
