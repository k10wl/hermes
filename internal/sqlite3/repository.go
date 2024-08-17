package sqlite3

import (
	"context"
	"database/sql"

	"github.com/k10wl/hermes/internal/models"
)

type executorFnSingle func(context.Context, string, ...interface{}) *sql.Row
type executorFnMultiple func(context.Context, string, ...interface{}) (*sql.Rows, error)

const createMessageQuery = `
INSERT INTO messages (chat_id, role_id, content)
VALUES ($1,$2,$3)
RETURNING id, chat_id, role_id, content, created_at, updated_at, deleted_at;
`

func createMessage(
	executor executorFnSingle,
	ctx context.Context,
	chatId int64,
	roleId int64,
	content string,
) (*models.Message, error) {
	row := executor(
		ctx,
		createMessageQuery,
		chatId,
		roleId,
		content,
	)
	var message models.Message
	err := row.Scan(
		&message.ID,
		&message.ChatID,
		&message.RoleID,
		&message.Content,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.DeletedAt,
	)
	return &message, err
}

const createChatQuery = `
INSERT INTO chats (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at, deleted_at;
`

func createChat(
	executor executorFnSingle,
	ctx context.Context,
	name string,
) (*models.Chat, error) {
	row := executor(ctx, createChatQuery, name)
	var chat models.Chat
	err := row.Scan(
		&chat.ID,
		&chat.Name,
		&chat.CreatedAt,
		&chat.UpdatedAt,
		&chat.DeletedAt,
	)
	return &chat, err
}

const getChatsQuery = `
SELECT id, name, created_at, updated_at, deleted_at FROM chats
ORDER BY created_at DESC;
`

func getChats(
	executor executorFnMultiple,
	ctx context.Context,
) ([]*models.Chat, error) {
	rows, err := executor(ctx, getChatsQuery)
	chats := []*models.Chat{}
	if err != nil {
		return chats, err
	}
	for rows.Next() {
		var chat models.Chat
		err = rows.Scan(
			&chat.ID,
			&chat.Name,
			&chat.CreatedAt,
			&chat.UpdatedAt,
			&chat.DeletedAt,
		)
		chats = append(chats, &chat)
		if err != nil {
			break
		}
	}
	return chats, err
}

const getChatMessagesQuery = `
SELECT id, chat_id, role_id, content, created_at, updated_at, deleted_at FROM messages
WHERE chat_id = $1;
`

func getChatMessages(
	executor executorFnMultiple,
	ctx context.Context,
	chatId int64,
) ([]*models.Message, error) {
	rows, err := executor(ctx, getChatMessagesQuery, chatId)
	messages := []*models.Message{}
	if err != nil {
		return messages, err
	}
	for rows.Next() {
		var message models.Message
		err = rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.RoleID,
			&message.Content,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.DeletedAt,
		)
		messages = append(messages, &message)
		if err != nil {
			break
		}
	}
	return messages, err
}

const getWebSettingsQuery = `
SELECT dark_mode, initted FROM web_settings;
`

func getWebSettings(
	executor executorFnSingle,
	ctx context.Context,
) (*models.WebSettings, error) {
	row := executor(ctx, getWebSettingsQuery)
	var webSettings models.WebSettings
	err := row.Scan(&webSettings.DarkMode, &webSettings.Initted)
	return &webSettings, err
}

const updateWebSettingsQuery = `
UPDATE web_settings 
SET initted = true, dark_mode = $1;
`

func updateWebSettings(executor executorFnSingle, ctx context.Context, darkMode bool) error {
	row := executor(ctx, updateWebSettingsQuery, darkMode)
	return row.Err()
}

const getLatestChatQuery = `
SELECT id, name, created_at, updated_at, deleted_at FROM chats
ORDER BY created_at DESC
LIMIT 1;
`

func getLatestChat(executor executorFnSingle, ctx context.Context) (*models.Chat, error) {
	row := executor(ctx, getLatestChatQuery)
	var chat models.Chat
	err := row.Scan(
		&chat.ID,
		&chat.Name,
		&chat.CreatedAt,
		&chat.UpdatedAt,
		&chat.DeletedAt,
	)
	return &chat, err
}
