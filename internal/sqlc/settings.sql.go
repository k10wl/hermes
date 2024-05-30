// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: settings.sql

package sqlc

import (
	"context"
)

const getSettings = `-- name: GetSettings :one
SELECT dark_mode FROM web_settings LIMIT 1
`

func (q *Queries) GetSettings(ctx context.Context) (bool, error) {
	row := q.db.QueryRowContext(ctx, getSettings)
	var dark_mode bool
	err := row.Scan(&dark_mode)
	return dark_mode, err
}

const updateSettings = `-- name: UpdateSettings :exec
UPDATE web_settings SET dark_mode = ?
`

func (q *Queries) UpdateSettings(ctx context.Context, darkMode bool) error {
	_, err := q.db.ExecContext(ctx, updateSettings, darkMode)
	return err
}