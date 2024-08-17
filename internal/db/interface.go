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
		roleId int64,
		content string,
	) (*models.Message, error)
	CreateChatAndMessage(
		ctx context.Context,
		roleId int64,
		content string,
	) (*models.Chat, *models.Message, error)
	GetChats(context.Context) ([]*models.Chat, error)
	GetChatMessages(context.Context, int64) ([]*models.Message, error)

	GetWebSettings(context.Context) (*models.WebSettings, error)
	UpdateWebSettings(ctx context.Context, dark_mode bool) error

	GetLatestChat(context.Context) (*models.Chat, error)
}
