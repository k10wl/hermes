package core

import (
	"context"

	"github.com/k10wl/hermes/internal/models"
)

type Query interface {
	Execute() error
}

type GetChatsQuery struct {
	Core   *Core
	Result []*models.Chat
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
	Result []*models.Message
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
	Result *models.WebSettings
}

func (q *WebSettingsQuery) Execute(ctx context.Context) error {
	if err := q.Core.assertAI(); err != nil {
		return err
	}
	setting, err := q.Core.db.GetWebSettings(ctx)
	q.Result = setting
	return err
}

type LatestChatQuery struct {
	Core   *Core
	Result *models.Chat
}

func (q *LatestChatQuery) Execute(ctx context.Context) error {
	chat, err := q.Core.db.GetLatestChat(ctx)
	q.Result = chat
	return err
}

type GetTemplateByNameQuery struct {
	Core   *Core
	Result *models.Template
	name   string
}

func NewGetTemplateByNameQuery(c *Core, name string) *GetTemplateByNameQuery {
	return &GetTemplateByNameQuery{Core: c, name: name}
}

func (q *GetTemplateByNameQuery) Execute(ctx context.Context) error {
	template, err := q.Core.db.GetTemplateByName(ctx, q.name)
	q.Result = template
	return err
}
