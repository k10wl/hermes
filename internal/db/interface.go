package db

import (
	"context"

	"github.com/k10wl/hermes/internal/models"
)

type Client interface {
	Close() error

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
	GetChats(
		ctx context.Context,
		limit int64,
		startBeforeID int64,
	) ([]*models.Chat, error)
	GetChatMessages(
		ctx context.Context,
		chatID int64,
	) ([]*models.Message, error)

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
	GetTemplates(
		ctx context.Context,
		after int64,
		limit int64,
		name string,
	) ([]*models.Template, error)
	DeleteTemplateByName(
		ctx context.Context,
		name string,
	) (bool, error)
	EditTemplateByName(
		ctx context.Context,
		name string,
		content string,
	) (bool, error)

	CreateActiveSession(*models.ActiveSession) error
	RemoveActiveSession(*models.ActiveSession) error
	GetActiveSessionByDatabaseDNS(string) (*models.ActiveSession, error)
}
