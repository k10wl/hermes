package sqlite3

import (
	"context"

	"github.com/k10wl/hermes/internal/models"
)

func (s *SQLite3) CreateMessage(
	ctx context.Context,
	chatId int64,
	role string,
	content string,
) (*models.Message, error) {
	return createMessage(s.DB.QueryRowContext, ctx, chatId, role, content)
}

func (s *SQLite3) CreateChat(ctx context.Context, name string) (*models.Chat, error) {
	return createChat(s.DB.QueryRowContext, ctx, name)
}

func (s *SQLite3) CreateChatAndMessage(
	ctx context.Context,
	role string,
	content string,
) (*models.Chat, *models.Message, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()
	chat, err := createChat(tx.QueryRowContext, ctx, content)
	if err != nil {
		return nil, nil, err
	}
	message, err := createMessage(tx.QueryRowContext, ctx, chat.ID, role, content)
	if err != nil {
		return nil, nil, err
	}
	err = tx.Commit()
	return chat, message, err
}

func (s *SQLite3) GetChats(
	ctx context.Context,
	limit int64,
	startAfterID int64,
) ([]*models.Chat, error) {
	return getChats(s.DB.QueryContext, ctx, limit, startAfterID)
}

func (s *SQLite3) GetChatMessages(
	ctx context.Context,
	chatID int64,
) ([]*models.Message, error) {
	return getChatMessages(s.DB.QueryContext, ctx, chatID)
}

func (s *SQLite3) GetWebSettings(ctx context.Context) (*models.WebSettings, error) {
	return getWebSettings(s.DB.QueryRowContext, ctx)
}

func (s *SQLite3) UpdateWebSettings(
	ctx context.Context,
	darkMode bool,
) error {
	return updateWebSettings(s.DB.QueryRowContext, ctx, darkMode)
}

func (s *SQLite3) GetLatestChat(
	ctx context.Context,
) (*models.Chat, error) {
	return getLatestChat(s.DB.QueryRowContext, ctx)
}

func (s SQLite3) UpsertTemplate(
	ctx context.Context,
	name string,
	template string,
) (*models.Template, error) {
	return upsertTemplate(s.DB.QueryRowContext, ctx, name, template)
}

func (s SQLite3) GetTemplatesByNames(
	ctx context.Context,
	name []string,
) ([]*models.Template, error) {
	return getTemplatesByNames(s.DB.QueryContext, ctx, name)
}

func (s SQLite3) DeleteTemplateByName(
	ctx context.Context,
	name string,
) (bool, error) {
	return deleteTemplateByName(s.DB.ExecContext, ctx, name)
}

func (s SQLite3) EditTemplateByName(
	ctx context.Context,
	name string,
	newName string,
	content string,
) (*models.Template, error) {
	return editTemplateByName(s.DB.QueryRowContext, ctx, name, newName, content)
}

func (s SQLite3) CreateActiveSession(activeSession *models.ActiveSession) error {
	return createActiveSession(s.DB.ExecContext, context.Background(), activeSession)
}

func (s SQLite3) RemoveActiveSession(activeSession *models.ActiveSession) error {
	return removeActiveSession(s.DB.ExecContext, context.Background(), activeSession)
}

func (s SQLite3) GetActiveSessionByDatabaseDNS(databaseDNS string) (*models.ActiveSession, error) {
	return getActiveSession(s.DB.QueryRowContext, context.Background(), databaseDNS)
}

func (s SQLite3) GetTemplates(
	ctx context.Context,
	after int64,
	limit int64,
	name string,
) ([]*models.Template, error) {
	return getTemplates(s.DB.QueryContext, ctx, after, limit, name)
}

func (s SQLite3) GetTemplateByID(
	ctx context.Context,
	id int64,
) (*models.Template, error) {
	return getTemplateByID(s.DB.QueryRowContext, ctx, id)
}
