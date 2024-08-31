package db

import (
	"context"

	"github.com/k10wl/hermes/internal/models"
)

type Client interface {
	CreateChat(context.Context, string) (*models.Chat, error)
	CreateMessage(
		ctx context.Context,
		chatId int64,
		role string,
		content string,
	) (*models.Message, error)
	CreateChatAndMessage(
		ctx context.Context,
		role string,
		content string,
	) (*models.Chat, *models.Message, error)
	GetChats(context.Context) ([]*models.Chat, error)
	GetChatMessages(context.Context, int64) ([]*models.Message, error)

	GetWebSettings(context.Context) (*models.WebSettings, error)
	UpdateWebSettings(ctx context.Context, dark_mode bool) error

	GetLatestChat(context.Context) (*models.Chat, error)

	UpsertTemplate(
		ctx context.Context,
		name string,
		template string,
	) (*models.Template, error)
	GetTemplatesByNames(
		ctx context.Context,
		names []string,
	) ([]*models.Template, error)
	GetTemplatesByRegexp(
		ctx context.Context,
		regexp string,
	) ([]*models.Template, error)
	DeleteTemplateByName(
		ctx context.Context,
		name string,
	) (bool, error)
}
