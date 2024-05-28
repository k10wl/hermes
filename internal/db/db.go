package db

import (
	"context"

	"github.com/k10wl/hermes/internal/sqlc"
)

type Client interface {
	CreateChat(context.Context, string) (sqlc.Chat, error)
	CreateMessage(context.Context, sqlc.CreateMessageParams) (sqlc.Message, error)
	CreateChatAndMessage(context.Context, sqlc.CreateMessageParams) (sqlc.Chat, sqlc.Message, error)
	GetChats(context.Context) ([]sqlc.Chat, error)
	GetChatMessages(context.Context, int64) ([]sqlc.GetChatMessagesRow, error)
}
