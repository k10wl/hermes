package v1

import (
	"net/http"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
)

func AddRoutes(mux *http.ServeMux, core *core.Core, hub *Hub) {
	mux.Handle("/api/v1/chats", handleChats(core))
	mux.Handle("/api/v1/health-check", handleCheckHeath())
	mux.Handle("/api/v1/relay", handleRelay(hub.broadcast))
	mux.Handle("/api/v1/ws", handleServeWebSockets(core, hub, ai_clients.Complete))
}
