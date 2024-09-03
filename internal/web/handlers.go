package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
)

func handleChat(c *core.Core, t *template.Template) http.HandlerFunc {
	type home struct {
		Chats       []*models.Chat
		Messages    []*models.Message
		WebSettings string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		data := home{}
		getChats := core.GetChatsQuery{
			Core: c,
		}
		err := getChats.Execute(context.Background())
		if err != nil {
			panic(err)
		}
		getSettings := core.WebSettingsQuery{Core: c}
		err = getSettings.Execute(context.Background())
		if err != nil {
			panic(err)
		}
		s, err := json.Marshal(getSettings.Result)
		data.WebSettings = string(s)
		data.Chats = getChats.Result
		chatId := r.PathValue("id")
		if id, err := strconv.ParseInt(chatId, 10, 64); err == nil {
			getMessages := core.GetChatMessagesQuery{Core: c, ChatID: id}
			getMessages.Execute(context.Background())
			data.Messages = getMessages.Result
		}
		t.ExecuteTemplate(w, "/home", data)
	}
}

func handleAssets() http.Handler {
	subFS, err := fs.Sub(assetsEmbed, "assets")
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.FS(subFS))
	return http.StripPrefix("/assets/", fs)
}

func handleMessage(c *core.Core, t *template.Template) http.HandlerFunc {
	type message struct {
		Content string
		Role    string
		ID      int64
		ChatID  int64
	}
	return func(w http.ResponseWriter, r *http.Request) {
		m := message{}
		content := r.FormValue("content")
		chatId := r.PathValue("id")
		// TODO handle error
		if chatId == "" {
			command := core.NewCreateChatAndCompletionCommand(
				c,
				core.UserRole,
				content,
				"",
				&ai_clients.Parameters{
					Model: "gpt-4o-mini",
				},
				ai_clients.Complete,
			)
			command.Execute(r.Context())
			m.Content = command.Result.Content
			m.ChatID = command.Result.ChatID
			m.Role = command.Result.Role
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Eval", "js")
			w.WriteHeader(http.StatusMovedPermanently)
			w.Write(
				[]byte(
					fmt.Sprintf(
						`window.location.replace('/chats/%v');`,
						m.ChatID,
					),
				),
			)
			return
		} else {
			id, err := strconv.ParseInt(chatId, 10, 64)
			if err != nil {
				panic(err)
			}
			command := core.NewCreateCompletionCommand(
				c,
				id,
				core.UserRole,
				content,
				"",
				&ai_clients.Parameters{
					Model: "gpt-4o-mini",
				},
				ai_clients.Complete,
			)
			command.Execute(context.Background())
			m.Content = command.Result.Content
			m.ID = command.Result.ID
			m.Role = command.Result.Role
		}
		t.ExecuteTemplate(w, "message", m)
	}
}

func handlePutSettings(c *core.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var s models.WebSettings
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &s)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		updateWebSettings := core.NewUpdateWebSettingsCommand(c, s)
		err = updateWebSettings.Execute(r.Context())
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
