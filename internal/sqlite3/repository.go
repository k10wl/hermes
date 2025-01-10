package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/k10wl/hermes/internal/models"
)

type queryRow func(context.Context, string, ...interface{}) *sql.Row
type queryRows func(context.Context, string, ...interface{}) (*sql.Rows, error)
type execute func(context.Context, string, ...interface{}) (sql.Result, error)

const createMessageQuery = `
INSERT INTO messages (chat_id, role_id, content)
VALUES ($1,$2,$3)
RETURNING id, chat_id, content, created_at, updated_at, deleted_at;
`

const getRoleIDByName = `
SELECT id FROM roles WHERE name = $1;
`

func createMessage(
	executor queryRow,
	ctx context.Context,
	chatId int64,
	role string,
	content string,
) (*models.Message, error) {
	row := executor(ctx, getRoleIDByName, role)
	var roleID int64
	err := row.Scan(&roleID)
	if err != nil {
		return nil, err
	}
	row = executor(
		ctx,
		createMessageQuery,
		chatId,
		roleID,
		content,
	)
	var message models.Message
	message.Role = role
	err = row.Scan(
		&message.ID,
		&message.ChatID,
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
	executor queryRow,
	ctx context.Context,
	name string,
) (*models.Chat, error) {
	row := executor(ctx, createChatQuery, ellipsis(name, 80, 3, "."))
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

func ellipsis(text string, max int, times int, replacement string) string {
	if len(text) <= max {
		return text
	}
	cutLen := len(replacement) * times
	if max-cutLen < 1 {
		panic("cannot trim into nothing")
	}
	cut := text[:max-cutLen]
	return fmt.Sprintf("%s%s", cut, strings.Repeat(replacement, times))
}

const getChatsQueryWithWhere = `
SELECT 
    id, name, created_at, updated_at, deleted_at
FROM chats
WHERE id < ?
ORDER BY id DESC
LIMIT ?;
`

const getChatsQuery = `
SELECT 
    id, name, created_at, updated_at, deleted_at
FROM chats
ORDER BY id DESC
LIMIT ?;
`

func getChats(
	executor queryRows,
	ctx context.Context,
	limit int64,
	startBeforeID int64,
) ([]*models.Chat, error) {
	var rows *sql.Rows
	var err error
	if startBeforeID < 1 {
		rows, err = executor(ctx, getChatsQuery, limit)
	} else {
		rows, err = executor(ctx, getChatsQueryWithWhere, startBeforeID, limit)
	}
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
SELECT m.id, chat_id, roles.name, content, m.created_at, m.updated_at, m.deleted_at 
FROM messages AS m
JOIN roles ON role_id = roles.id
WHERE chat_id = $1;
`

func getChatMessages(
	executor queryRows,
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
			&message.Role,
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
	executor queryRow,
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

func updateWebSettings(executor queryRow, ctx context.Context, darkMode bool) error {
	row := executor(ctx, updateWebSettingsQuery, darkMode)
	return row.Err()
}

const getLatestChatQuery = `
SELECT id, name, created_at, updated_at, deleted_at FROM chats
ORDER BY id DESC
LIMIT 1;
`

func getLatestChat(executor queryRow, ctx context.Context) (*models.Chat, error) {
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

const upsertTemplateQuery = `
INSERT OR REPLACE INTO templates (name, content)
VALUES ($1, $2)
RETURNING id, name, content, created_at, updated_at, deleted_at;
`

func upsertTemplate(
	executor queryRow,
	ctx context.Context,
	name string,
	content string,
) (*models.Template, error) {
	row := executor(ctx, upsertTemplateQuery, name, content)
	var templateDoc models.Template
	err := row.Scan(
		&templateDoc.ID,
		&templateDoc.Name,
		&templateDoc.Content,
		&templateDoc.CreatedAt,
		&templateDoc.UpdatedAt,
		&templateDoc.DeletedAt,
	)
	return &templateDoc, err
}

func getTemplatesByNamesQuery(names []interface{}) string {
	// my brother in christ this is painful to write... holy fuck
	return `
SELECT id, name, content, created_at, updated_at, deleted_at FROM templates
WHERE name IN (?` + strings.Repeat(",?", len(names)-1) + `) AND deleted_at IS NULL;`
}

func getTemplatesByNames(
	executor queryRows,
	ctx context.Context,
	names []string,
) ([]*models.Template, error) {
	namesInterface := convertToAnySlice(names)
	rows, err := executor(ctx, getTemplatesByNamesQuery(namesInterface), namesInterface...)
	if err != nil {
		return nil, err
	}
	templates := []*models.Template{}
	var rowErr error
	for rows.Next() {
		var templateDoc models.Template
		if err := rows.Scan(
			&templateDoc.ID,
			&templateDoc.Name,
			&templateDoc.Content,
			&templateDoc.CreatedAt,
			&templateDoc.UpdatedAt,
			&templateDoc.DeletedAt,
		); err != nil {
			rowErr = err
			break
		}
		templates = append(templates, &templateDoc)
	}
	return templates, rowErr
}

const getTemplatesQueryBase = `
SELECT
    id,
    name,
    content,
    created_at,
    updated_at,
    deleted_at
FROM
    templates`

var getTemplatesAfterQuery = fmt.Sprintf(`
%s
WHERE
    deleted_at IS NULL AND
    id < ? AND
    name LIKE ?
ORDER BY id DESC
LIMIT ?;`, getTemplatesQueryBase)

var getTemplatesQuery = fmt.Sprintf(`
%s
WHERE
    deleted_at IS NULL AND
    name LIKE ?
ORDER BY id DESC
LIMIT ?;`, getTemplatesQueryBase)

func scanTemplate(scan func(dest ...any) error, receiver *models.Template) error {
	return scan(
		&receiver.ID,
		&receiver.Name,
		&receiver.Content,
		&receiver.CreatedAt,
		&receiver.UpdatedAt,
		&receiver.DeletedAt,
	)
}

func getTemplates(
	executor queryRows,
	ctx context.Context,
	after int64,
	limit int64,
	name string,
) ([]*models.Template, error) {
	refinedName := "%" + name + "%"
	var rows *sql.Rows
	var err error
	if after > 0 {
		rows, err = executor(ctx, getTemplatesAfterQuery, after, refinedName, limit)
	} else {
		rows, err = executor(ctx, getTemplatesQuery, refinedName, limit)
	}
	if err != nil {
		return nil, err
	}
	templates := []*models.Template{}
	var rowErr error
	for rows.Next() {
		var templateDoc models.Template
		if err := scanTemplate(rows.Scan, &templateDoc); err != nil {
			rowErr = err
			break
		}
		templates = append(templates, &templateDoc)
	}
	return templates, rowErr
}

var getTemplateByIDQuery = fmt.Sprintf(`%s
WHERE
    id = ?`, getTemplatesQueryBase)

func getTemplateByID(
	executor queryRow,
	ctx context.Context,
	id int64,
) (*models.Template, error) {
	row := executor(ctx, getTemplateByIDQuery, id)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	var template models.Template
	return &template, scanTemplate(row.Scan, &template)
}

const deleteTemplateByNameQuery = `
UPDATE templates
SET deleted_at = CURRENT_TIMESTAMP
WHERE name = $1;`

func deleteTemplateByName(
	executor execute,
	ctx context.Context,
	name string,
) (bool, error) {
	res, err := executor(ctx, deleteTemplateByNameQuery, name)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	return affected == 1, err
}

const editTemplateByNameQuery = `
UPDATE templates
SET 
    content = ?,
    updated_at = ?
WHERE name = ?
RETURNING id, name, content, created_at, updated_at, deleted_at
`

func editTemplateByName(
	executor queryRow,
	ctx context.Context,
	name string,
	content string,
) (*models.Template, error) {
	result := executor(ctx, editTemplateByNameQuery, content, time.Now(), name)
	if err := result.Err(); err != nil {
		return nil, err
	}
	tmp := new(models.Template)
	if err := scanTemplate(result.Scan, tmp); err != nil {
		return nil, err
	}
	return tmp, nil
}

const createActiveSessionQuery = `
INSERT INTO active_sessions (address, database_dns)
VALUES ($1, $2);
`

func createActiveSession(
	executor execute,
	ctx context.Context,
	activeSession *models.ActiveSession,
) error {
	res, err := executor(
		ctx,
		createActiveSessionQuery,
		activeSession.Address,
		activeSession.DatabaseDNS,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("failed to create active session\n")
	}
	return err
}

const removeActiveSessionQuery = `
DELETE FROM active_sessions WHERE (database_dns = $1);
`

func removeActiveSession(
	executor execute,
	ctx context.Context,
	activeSession *models.ActiveSession,
) error {
	res, err := executor(
		ctx,
		removeActiveSessionQuery,
		activeSession.DatabaseDNS,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("failed to remove active session\n")
	}
	return err
}

const getActiveSessionsQuery = `
SELECT id, address, database_dns FROM active_sessions
WHERE database_dns = $1;
`

func getActiveSession(
	executor queryRow,
	ctx context.Context,
	databaseDNS string,
) (*models.ActiveSession, error) {
	res := executor(ctx, getActiveSessionsQuery, databaseDNS)
	err := res.Err()
	if err != nil {
		return nil, err
	}
	activeSession := models.ActiveSession{}
	err = res.Scan(&activeSession.ID, &activeSession.Address, &activeSession.DatabaseDNS)
	if err != nil {
		return nil, err
	}
	return &activeSession, nil
}
