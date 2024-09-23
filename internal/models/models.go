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
	MaxTokens   int64 `json:"max_tokens"`
	Temperature int64 `json:"temperature"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Model string `json:"model"`
	Name  string `json:"name"`
	CompletionParameters
	Timestamps
}

type Message struct {
	ID                 int64  `json:"id"`
	ChatID             int64  `json:"chat_id"`
	Role               string `json:"role"`
	Content            string `json:"content"`
	Generation         int64  `json:"generation"`
	SelectedGeneration bool   `json:"selected_generation"`
	Timestamps
}

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type WebSettings struct {
	DarkMode bool `json:"dark_mode"`
	Initted  bool `json:"initted"`
}

type Template struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Timestamps
}

type ActiveSession struct {
	ID          int64  `json:"id"`
	Address     string `json:"address"`
	DatabaseDNS string `json:"database_dns"`
}

// if you used this anywhere (any-fucking-where) outside of test function - you fucked up
func (t *Timestamps) TimestampsToNilForTest__() {
	t.CreatedAt = nil
	t.UpdatedAt = nil
	t.DeletedAt = nil
}
