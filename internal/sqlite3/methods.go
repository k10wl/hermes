package sqlite3

import (
	"context"

	"github.com/k10wl/hermes/internal/models"
)

func (s *SQLite3) CreateMessage(
	ctx context.Context,
	chatId int64,
	roleId int64,
	content string,
) (*models.Message, error) {
	return createMessage(s.db.QueryRowContext, ctx, chatId, roleId, content)
}

func (s *SQLite3) CreateChat(ctx context.Context, name string) (*models.Chat, error) {
	return createChat(s.db.QueryRowContext, ctx, name)
}

func (s *SQLite3) CreateChatAndMessage(
	ctx context.Context,
	roleId int64,
	content string,
) (*models.Chat, *models.Message, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()
	chat, err := createChat(tx.QueryRowContext, ctx, content)
	if err != nil {
		return nil, nil, err
	}
	message, err := createMessage(tx.QueryRowContext, ctx, chat.ID, 1, content)
	if err != nil {
		return nil, nil, err
	}
	err = tx.Commit()
	return chat, message, err
}

func (s *SQLite3) GetChats(ctx context.Context) ([]*models.Chat, error) {
	return getChats(s.db.QueryContext, ctx)
}

func (s *SQLite3) GetChatMessages(
	ctx context.Context,
	chatID int64,
) ([]*models.Message, error) {
	return getChatMessages(s.db.QueryContext, ctx, chatID)
}

func (s *SQLite3) GetWebSettings(ctx context.Context) (*models.WebSettings, error) {
	return getWebSettings(s.db.QueryRowContext, ctx)
}

func (s *SQLite3) UpdateWebSettings(
	ctx context.Context,
	darkMode bool,
) error {
	return updateWebSettings(s.db.QueryRowContext, ctx, darkMode)
}

func (s *SQLite3) GetLatestChat(
	ctx context.Context,
) (*models.Chat, error) {
	return getLatestChat(s.db.QueryRowContext, ctx)
}
