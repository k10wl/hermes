package v1

import (
	"net/http"

	"github.com/k10wl/hermes/internal/core"
)

func AddRoutes(mux *http.ServeMux, core *core.Core) {
	mux.Handle("/api/v1/chats", handleChats(core))
}
