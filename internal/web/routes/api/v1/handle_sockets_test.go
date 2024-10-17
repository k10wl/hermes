package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/k10wl/hermes/internal/test_helpers"
)

func TestHandleWebSocket(t *testing.T) {
	c, _ := test_helpers.CreateCore()
	hub := NewHub()
	go hub.Run()
	server := httptest.NewServer(http.HandlerFunc(handleServeWebSockets(c, hub)))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		t.Fatalf("err reading res %s\n", err)
	}
	if err != nil {
		t.Fatalf("could not connect to WebSocket server: %v", err)
	}
	defer client.Close()

	message, _ := newMessage("hello", nil).encode()
	err = client.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	_, response, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	if string(response) != string(message) {
		t.Errorf("expected response '%s', but got '%s'", string(message), string(response))
	}
}
