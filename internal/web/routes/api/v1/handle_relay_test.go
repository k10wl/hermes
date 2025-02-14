package v1

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
)

func TestRelay(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	coreInstance, database := test_helpers.CreateCore()
	mux := http.NewServeMux()
	AddRoutes(mux, coreInstance, hub)
	server := httptest.NewServer(mux)
	defer server.Close()
	ctx := context.Background()
	conn, teardownConn, err := test_helpers.CreateWebsocketConnection(
		"ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws",
	)
	defer teardownConn()
	if err != nil {
		t.Fatalf("failed to create websocket connection - %s\n", err)
	}

	seeder := db_helpers.NewSeeder(database, ctx)
	if err := seeder.SeedChatsN(1); err != nil {
		t.Fatalf("failed to seed chats - %s\n", err)
	}
	if err := seeder.SeedMessagesN(2, 1); err != nil {
		t.Fatalf("failed to seed messages - %s\n", err)
	}
	dbChat, err := db_helpers.GetChatByID(database, ctx, 1)
	if err != nil {
		t.Fatalf("failed to seed messages - %s\n", err)
	}
	dbMessages, err := db_helpers.GetMessagesByChatID(database, ctx, 1)
	if err != nil {
		t.Fatalf("failed to seed messages - %s\n", err)
	}

	firstMessage, err := messages.Encode(
		messages.NewServerChatCreated(uuid.NewString(), dbChat, dbMessages[0]),
	)
	if err != nil {
		t.Fatalf("failed to encode message - %s\n", err)
	}
	res, err := http.Post(
		server.URL+"/api/v1/relay",
		"text/plain",
		bytes.NewReader(firstMessage),
	)
	if err != nil {
		t.Fatalf("failed to post into relay - %s\n", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf(
			"failed to get successful status code, actual: %d\n",
			res.StatusCode,
		)
	}
	if err := conn.ReadJSON(&messages.ServerChatCreated{}); err != nil {
		t.Fatalf("failed to read websocket message - %s\n", err)
	}

	secondMessage, err := messages.Encode(
		messages.NewServerMessageCreated(
			uuid.NewString(),
			dbChat.ID,
			dbMessages[1],
		),
	)
	if err != nil {
		t.Fatalf("failed to encode second message - %s\n", err)
	}
	res, err = http.Post(
		server.URL+"/api/v1/relay",
		"text/plain",
		bytes.NewReader(secondMessage),
	)
	if err != nil {
		t.Fatalf("failed to post second message into relay - %s\n", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf(
			"failed to get second successful status code, actual: %d\n",
			res.StatusCode,
		)
	}
	if err := conn.ReadJSON(&messages.ServerMessageCreated{}); err != nil {
		t.Fatalf("failed to read second websocket message - %s\n", err)
	}
}
