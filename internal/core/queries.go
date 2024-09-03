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

type GetTemplatesByNamesQuery struct {
	Core   *Core
	Result []*models.Template
	names  []string
}

func NewGetTemplatesByNamesQuery(c *Core, names []string) *GetTemplatesByNamesQuery {
	return &GetTemplatesByNamesQuery{Core: c, names: names}
}

func (q *GetTemplatesByNamesQuery) Execute(ctx context.Context) error {
	template, err := q.Core.db.GetTemplatesByNames(ctx, q.names)
	q.Result = template
	return err
}

type GetTemplatesByRegexp struct {
	Core   *Core
	Result []*models.Template
	regexp string
}

func NewGetTemplatesByRegexp(c *Core, regexp string) *GetTemplatesByRegexp {
	return &GetTemplatesByRegexp{Core: c, regexp: regexp}
}

func (q *GetTemplatesByRegexp) Execute(ctx context.Context) error {
	templates, err := q.Core.db.GetTemplatesByRegexp(ctx, q.regexp)
	q.Result = templates
	return err
}
