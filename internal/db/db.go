package db

import (
	"context"
	"database/sql"

	"github.com/k10wl/hermes/internal/sqlc"
)

type Client interface {
	CreateChat(context.Context, sql.NullString) (sqlc.Chat, error)
	CreateMessage(context.Context, sqlc.CreateMessageParams) (sqlc.Message, error)
	CreateChatAndMessage(context.Context, sqlc.CreateMessageParams) (sqlc.Chat, sqlc.Message, error)
}
