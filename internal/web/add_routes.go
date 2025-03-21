package web

import (
	"html/template"
	"net/http"

	"github.com/k10wl/hermes/internal/core"
	v1 "github.com/k10wl/hermes/internal/web/routes/api/v1"
)

func addRoutes(mux *http.ServeMux, core *core.Core, hub *v1.Hub, t *template.Template) {
	mux.Handle("/", handleChat(core, t))
	mux.Handle("/chats/{id}", handleChat(core, t))
	mux.Handle("POST /chats", handleMessage(core, t))
	mux.Handle("POST /chats/{id}", handleMessage(core, t))
	mux.Handle("/assets/", handleAssets())
	mux.Handle("PUT /settings", handlePutSettings(core))

	v1.AddRoutes(mux, core, hub)
}
