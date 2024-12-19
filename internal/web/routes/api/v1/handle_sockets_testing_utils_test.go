package v1

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/k10wl/hermes/internal/test_helpers"
)

const sharedID = "717dc403-63ab-48e6-94e8-21b3110da18c"

func setupWebSocketTest(t *testing.T) (*websocket.Conn, *sql.DB, func()) {
	c, db := test_helpers.CreateCore()
	hub := NewHub()
	go hub.Run()
	server := httptest.NewServer(
		http.HandlerFunc(
			handleServeWebSockets(c, hub, test_helpers.MockCompletion),
		),
	)
	client, cleanupClient, err := test_helpers.CreateWebsocketConnection(
		"ws" + strings.TrimPrefix(server.URL, "http"),
	)
	if err != nil {
		t.Fatalf("could not connect to WebSocket server: %v", err)
	}
	return client, db, func() {
		cleanupClient()
		client.Close()
		server.Close()
		db.Close()
	}
}
