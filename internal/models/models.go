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
	ID      int64  `json:"id"      validate:"required"`
	Name    string `json:"name"    validate:"required"`
	Content string `json:"content" validate:"required"`
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

var DefaultTemplate = Template{
	ID:   -1,
	Name: "example name",
	Content: `--{{define "example name"}}
<Prompt>
This block defines an example upsert template.
Prompt XML tag is not required, it helps to ` + "`cat`" + ` this text in VIM

Quick info about template capabilities:
>>> --{{.}} - Prints entire template input. If empty - prints <no value>

>>> --{{with .}}
      --{{.}} - Will print temlate input only if it is not empty
  --{{end}}

>>> --{{if .isEnabled}}
      --{{.jsonKey}} - prints out specific json key
      This block runs if the condition is true.
  --{{end}}

>>> Usefull examles from [docs](https://pkg.go.dev/text/template#hdr-Actions)
</Prompt>
--{{end}}
`,
}
