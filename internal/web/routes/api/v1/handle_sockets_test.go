package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
)

func TestHandleWebSocketPing(t *testing.T) {
	client, _, teardown := setupWebSocketTest(t)
	defer teardown()

	err := client.WriteMessage(websocket.TextMessage, []byte(`{"type": "ping"}`))
	if err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	_, response, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	expected := `{"type":"pong"}`
	if string(response) != expected {
		t.Errorf("expected response '%s', but got '%s'", expected, string(response))
	}
}

func TestRequestReadChat(t *testing.T) {
	client, db, teardown := setupWebSocketTest(t)
	defer teardown()

	seeder := db_helpers.NewSeeder(db, context.TODO())
	err := seeder.SeedChatsN(1)
	if err != nil {
		t.Fatal(err)
	}
	messagesAmount := 5
	err = seeder.SeedMessagesN(int64(messagesAmount), 1)
	if err != nil {
		t.Fatal(err)
	}

	err = client.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"type": "request-read-chat", "payload": 1}`),
	)
	if err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	_, response, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	res := messages.ReadChatMessage{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		t.Errorf("Failed to decode server response message\n")
	}
	if res.Payload.ChatID != 1 {
		t.Errorf("Failed to return same chat\n")
	}
	if len(res.Payload.Messages) != messagesAmount {
		// NOTE might change if pagination will be needed
		t.Errorf(
			"Failed to return all messages from chat, expected %d, got %d\n",
			messagesAmount,
			len(res.Payload.Messages),
		)
	}
}

func setupWebSocketTest(t *testing.T) (*websocket.Conn, *sql.DB, func()) {
	c, db := test_helpers.CreateCore()
	hub := NewHub()
	go hub.Run()
	server := httptest.NewServer(http.HandlerFunc(handleServeWebSockets(c, hub)))
	url := "ws" + strings.TrimPrefix(server.URL, "http")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	client, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		t.Fatalf("could not connect to WebSocket server: %v", err)
	}
	return client, db, func() {
		cancel()
		client.Close()
		server.Close()
		db.Close()
	}
}
