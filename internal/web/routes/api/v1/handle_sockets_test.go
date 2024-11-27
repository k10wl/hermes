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

const sharedID = "717dc403-63ab-48e6-94e8-21b3110da18c"

func TestHandleWebSocketPing(t *testing.T) {
	client, _, teardown := setupWebSocketTest(t)
	defer teardown()

	err := client.WriteMessage(websocket.TextMessage, []byte(
		`{ "id": "717dc403-63ab-48e6-94e8-21b3110da18c", "type": "ping" }`,
	))
	if err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	_, response, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	expected := `{"id":"717dc403-63ab-48e6-94e8-21b3110da18c","type":"pong"}`
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
		[]byte(
			`{"type": "request-read-chat", "payload": 1, "id": "717dc403-63ab-48e6-94e8-21b3110da18c"}`,
		),
	)
	if err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	_, response, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	res := messages.ServerReadChat{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		t.Errorf("Failed to decode server response message - %s\n", response)
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

func TestCreateMessageInExistingChat(t *testing.T) {
	client, db, teardown := setupWebSocketTest(t)
	defer teardown()

	seeder := db_helpers.NewSeeder(db, context.TODO())
	err := seeder.SeedChatsN(1)
	if err != nil {
		t.Fatal(err)
	}

	err = client.WriteMessage(
		websocket.TextMessage,
		[]byte(`
{
  "id": "717dc403-63ab-48e6-94e8-21b3110da18c",
  "type": "create-completion",
  "payload": {
    "chat_id": 1,
    "content": "create message",
    "template": "",
    "parameters": {
      "model": "gpt-4o-mini"
    }
  }
}
`),
	)
	if err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	_, response, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	res := messages.ServerMessageCreated{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		t.Errorf("Failed to decode server response message - %s\n", response)
	}
	if res.Type != "message-created" {
		t.Errorf("Did not respond with 'message created'\n")
	}
	if res.ID != sharedID {
		t.Errorf("Failed to return shared id\n")
	}
	if res.Payload.ChatID != 1 {
		t.Errorf("Failed to return same chat\n")
	}
	if res.Payload.Message.Content != "create message" {
		t.Errorf("Failed to create message, got %q\n", res.Payload.Message.Content)
	}

	_, response, err = client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	res = messages.ServerMessageCreated{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		t.Errorf("Failed to decode server response message - %s\n", response)
	}
	if res.Type != "message-created" {
		t.Errorf("Did not respond with 'message created'\n")
	}
	if res.Payload.ChatID != 1 {
		t.Errorf("Failed to return same chat\n")
	}
	if res.ID != sharedID {
		t.Errorf("Failed to return shared id\n")
	}
	if res.Payload.Message.Content != "> mocked: create message" {
		t.Errorf("Failed to create message, got %q\n", res.Payload.Message.Content)
	}

	msg, err := db_helpers.GetMessagesByChatID(db, context.Background(), 1)
	if len(msg) != 2 {
		t.Fatalf(
			"Did not store correct amount of messages in database, expected 2, but got %d",
			len(msg),
		)
	}
}

func TestCreateMessageInNewChat(t *testing.T) {
	client, db, teardown := setupWebSocketTest(t)
	defer teardown()

	err := client.WriteMessage(
		websocket.TextMessage,
		[]byte(`
{
  "type": "create-completion",
  "id": "717dc403-63ab-48e6-94e8-21b3110da18c",
  "payload": {
    "chat_id": -1,
    "content": "create message",
    "template": "",
    "parameters": {
      "model": "gpt-4o-mini"
    }
  }
}
`),
	)
	if err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	_, response, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	res := messages.ServerChatCreated{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		t.Errorf("Failed to decode server response message - %s\n", response)
	}
	if res.Type != "chat-created" {
		t.Errorf("Did not respond with 'message created'\n")
	}
	if res.Payload.Chat.ID != 1 {
		t.Errorf("Failed to return new chat id\nChat: %+v\n", res.Payload.Chat)
	}
	if res.ID != sharedID {
		t.Errorf("Failed to return shared id\n")
	}
	if res.Payload.Message.Content != "create message" {
		t.Errorf("Failed to create desired message\nMessage: %+v\n", res.Payload.Message)
	}

	_, response, err = client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	res2 := messages.ServerMessageCreated{}
	err = json.Unmarshal(response, &res2)
	if err != nil {
		t.Errorf("Failed to decode server response message - %s\n", response)
	}
	if res2.Type != "message-created" {
		t.Errorf("Did not respond with 'message created'\n")
	}
	if res2.Payload.ChatID != 1 {
		t.Errorf("Failed to return new chat id\nPayload: %+v\n", res2.Payload)
	}
	if res.ID != sharedID {
		t.Errorf("Failed to return shared id\n")
	}
	if res2.Payload.Message.Content != "> mocked: create message" {
		t.Errorf("Failed to create desired message, Actual: %+v\n", res2.Payload.Message)
	}

	msg, err := db_helpers.GetMessagesByChatID(db, context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to retrieve stored messages, error: %s\n", err)
	}
	if len(msg) != 2 {
		t.Fatalf(
			"Did not store correct amount of messages in database, expected 2, but got %d",
			len(msg),
		)
	}
}

func TestCreateMessageInNonExistingChat(t *testing.T) {
	client, _, teardown := setupWebSocketTest(t)
	defer teardown()

	err := client.WriteMessage(
		websocket.TextMessage,
		[]byte(`
{
  "id": "717dc403-63ab-48e6-94e8-21b3110da18c",
  "type": "create-completion",
  "payload": {
    "chat_id": 999,
    "content": "create message",
    "template": "",
    "parameters": {
      "model": "gpt-4o-mini"
    }
  }
}
`),
	)
	if err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	_, response1, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	res1 := messages.ServerMessageCreated{}
	err = json.Unmarshal(response1, &res1)
	if err != nil {
		t.Errorf("Failed to decode server response message\nData: %s", response1)
	}
	if res1.Type != "message-created" {
		t.Errorf("Did not respond with 'message created'\n")
	}
	if res1.Payload.Message.Content != "create message" {
		t.Errorf("Failed to create desired message\nMessage: %+v\n", res1.Payload.Message)
	}
	if res1.ID != sharedID {
		t.Errorf("Failed to return shared id\n")
	}

	_, response2, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}
	res2 := messages.ServerError{}
	err = json.Unmarshal(response2, &res2)
	if err != nil {
		t.Errorf("Failed to decode server response message\nData: %s", response1)
	}
	if res2.Type != "server-error" {
		t.Errorf("Server did not throw expected error upon creating message in non existing chat\n")
	}
	if res2.ID != sharedID {
		t.Errorf("Failed to return shared id\n")
	}
}

func setupWebSocketTest(t *testing.T) (*websocket.Conn, *sql.DB, func()) {
	c, db := test_helpers.CreateCore()
	hub := NewHub()
	go hub.Run()
	server := httptest.NewServer(
		http.HandlerFunc(
			handleServeWebSockets(c, hub, test_helpers.MockCompletion),
		),
	)
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
