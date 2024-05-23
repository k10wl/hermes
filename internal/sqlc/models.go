// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlc

import (
	"database/sql"
	"time"
)

type Chat struct {
	ID        int64
	Name      sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type Message struct {
	ID        int64
	ChatID    int64
	Content   string
	RoleID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type Role struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
