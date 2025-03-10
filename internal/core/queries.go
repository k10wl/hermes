package core

import (
	"context"

	"github.com/k10wl/hermes/internal/models"
)

type Query interface {
	Execute(context.Context) error
}

type GetChatsQuery struct {
	core          *Core
	limit         int64
	startBeforeID int64
	Result        []*models.Chat
}

// limit -1 forces to return all results
func NewGetChatsQuery(core *Core, limit int64, startBeforeID int64) *GetChatsQuery {
	return &GetChatsQuery{
		core:          core,
		limit:         limit,
		startBeforeID: startBeforeID,
	}
}

func (q *GetChatsQuery) Execute(ctx context.Context) error {
	chats, err := q.core.db.GetChats(ctx, q.limit, q.startBeforeID)
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

type GetTemplatesQuery struct {
	core   *Core
	Result []*models.Template
	after  int64
	limit  int64
	name   string
}

func NewGetTemplatesQuery(
	c *Core,
	startBeforeID int64,
	limit int64,
	name string,
) *GetTemplatesQuery {
	return &GetTemplatesQuery{
		core:  c,
		after: startBeforeID,
		limit: limit,
		name:  name,
	}
}

func (q *GetTemplatesQuery) Execute(ctx context.Context) error {
	res, err := q.core.db.GetTemplates(ctx, q.after, q.limit, q.name)
	q.Result = res
	return err
}

type GetTemplateByIDQuery struct {
	core   *Core
	id     int64
	Result *models.Template
}

func NewGetTemplateByIDQuery(c *Core, id int64) *GetTemplateByIDQuery {
	return &GetTemplateByIDQuery{
		core: c,
		id:   id,
	}
}

func (q *GetTemplateByIDQuery) Execute(ctx context.Context) error {
	res, err := q.core.db.GetTemplateByID(ctx, q.id)
	q.Result = res
	return err
}
