package models

import (
	"database/sql"
	"time"
)

type Timestamps struct {
	CreatedAt *time.Time    `json:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at"`
	DeletedAt *sql.NullTime `json:"deleted_at"`
}

type CompletionParameters struct {
	MaxTokens   int64
	Temperature int64
}

type Model struct {
	ID        int64
	Provider  string
	Name      string
	MaxTokens int64
	Timestamps
}

type Chat struct {
	ID    int64
	Model *Model
	Name  string
	CompletionParameters
	Timestamps
}

type Message struct {
	ID                 int64  `json:"id"`
	ChatID             int64  `json:"chat_id"`
	RoleID             int64  `json:"role_id"`
	Content            string `json:"content"`
	Generation         int64  `json:"generation"`
	SelectedGeneration bool   `json:"selected_generation"`
	Timestamps
}

type Role struct {
	ID   int64
	Name string
}

type WebSettings struct {
	DarkMode bool
	Initted  bool
}
