package web

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/runtime"
	"github.com/k10wl/hermes/internal/sqlc"
)

//go:embed assets
var assetsEmbed embed.FS

//go:embed views
var viewsEmbed embed.FS

func Serve(core *core.Core, config *runtime.Config) error {
	server := NewServer(core)
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	httpServer := http.Server{
		Addr:    addr,
		Handler: server,
	}
	fmt.Printf("Starting server on %s\n", addr)
	return httpServer.ListenAndServe()
}

func NewServer(core *core.Core) http.Handler {
	mux := http.NewServeMux()
	t := NewTemplate()
	addRoutes(mux, core, t)
	return mux
}

func addRoutes(mux *http.ServeMux, core *core.Core, t *template.Template) {
	mux.Handle("/", handleChat(core, t))
	mux.Handle("/chats/{id}", handleChat(core, t))
	mux.Handle("POST /chats", handleMessage(core, t))
	mux.Handle("POST /chats/{id}", handleMessage(core, t))
	mux.Handle("/assets/", handleAssets())
}

func handleChat(c *core.Core, t *template.Template) http.HandlerFunc {
	type home struct {
		Chats    []sqlc.Chat
		Messages []sqlc.GetChatMessagesRow
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
			command := &core.CreateChatAndCompletionCommand{Core: c, Message: content}
			command.Execute(context.Background())
			m.Content = command.Result.Content
			m.ChatID = command.Result.ChatID
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Eval", "js")
			w.WriteHeader(http.StatusMovedPermanently)
			w.Write(
				[]byte(
					fmt.Sprintf(`window.location.replace('/chats/%v');`, m.ChatID),
				),
			)
			return
		} else {
			id, err := strconv.ParseInt(chatId, 10, 64)
			if err != nil {
				panic(err)
			}
			command := &core.CreateCompletionCommand{Core: c, ChatID: id, Message: content}
			command.Execute(context.Background())
			m.Content = command.Result.Content
			m.ID = command.Result.ID
		}
		m.Role = core.AssistantRole
		t.ExecuteTemplate(w, "message", m)
	}
}

func NewTemplate() *template.Template {
	tmpl := template.New("main")
	templateContent, err := viewsEmbed.ReadFile("views/home.html")
	if err != nil {
		panic(err)
	}
	tmpl, err = tmpl.Parse(string(templateContent))
	if err != nil {
		panic(err)
	}
	return tmpl
}
