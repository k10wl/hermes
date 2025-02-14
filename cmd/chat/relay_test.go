package chat_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/k10wl/hermes/cmd/chat"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
)

func TestRelayMessageWithNewChat(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	coreInstance.GetConfig().DatabaseDSN = "deez"
	cmd := chat.CreateChatCommand(coreInstance, test_helpers.MockCompletion)
	cmd.Flags().Set("content", "testing stuff")

	relayDataChan := make(chan string)
	requestIDChan := make(chan string)
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		id := r.Header.Get("ID")
		dets := fmt.Sprintf("%s | %s | %s", r.Method, r.URL.String(), body)
		go func() {
			relayDataChan <- dets
			requestIDChan <- id
		}()
	}))
	server := httptest.NewServer(mux)
	defer server.Close()

	db_helpers.CreateActiveSession(
		db,
		context.Background(),
		&models.ActiveSession{
			DatabaseDNS: coreInstance.GetConfig().DatabaseDSN,
			Address:     server.URL,
		},
	)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute cmd: %s", err)
	}

	relayData := []string{}
	idsData := []string{}
	for i := 0; i < 4; i++ {
		select {
		case data := <-relayDataChan:
			relayData = append(relayData, data)
		case data := <-requestIDChan:
			idsData = append(idsData, data)
		}
	}

	dbChat, err := db_helpers.GetChatByID(db, context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to retrieve db chat: %s", err)
	}
	dbMessages, err := db_helpers.GetMessagesByChatID(db, context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to retrieve db chat messages: %s", err)
	}

	serverChatCreatedMessage := messages.ServerChatCreated{
		ID:   idsData[0],
		Type: "chat-created",
		Payload: messages.ServerChatCreatedPayload{
			Chat:    dbChat,
			Message: dbMessages[0],
		},
	}

	chatCreatedData, err := messages.Encode(serverChatCreatedMessage)
	if err != nil {
		t.Fatalf("failed to encode expected chat message for test: %s", err)
	}

	expectedChatCreatedData := fmt.Sprintf("POST | /api/v1/relay | %s", chatCreatedData)
	if relayData[0] != expectedChatCreatedData {
		t.Fatalf(
			"failed to get chat create message event\nexpected: %v\nactual:   %v\n",
			expectedChatCreatedData,
			relayData[0],
		)
	}

	serverMessageCreated := messages.ServerMessageCreated{
		ID:   idsData[1],
		Type: "message-created",
		Payload: messages.ServerMessageCreatedPayload{
			ChatID:  1,
			Message: dbMessages[1],
		},
	}

	messageData, err := messages.Encode(serverMessageCreated)
	if err != nil {
		t.Fatalf("failed to encode expected chat message for test: %s", err)
	}

	expectedMessageData := fmt.Sprintf("POST | /api/v1/relay | %s", messageData)
	if relayData[1] != expectedMessageData {
		t.Fatalf(
			"failed to get chat create message event\nexpected: %v\nactual:   %v\n",
			expectedMessageData,
			relayData[1],
		)
	}

	out := coreInstance.GetConfig().Stdoout.(*strings.Builder)
	if out.String() != dbMessages[len(dbMessages)-1].Content+"\n" {
		t.Fatalf(
			"failed to write result to stdout\nexpected: %q\nactual:   %q\n",
			"> mocked: testing stuff\n",
			out.String(),
		)
	}
}

func TestRelayMessageWithExistingChat(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	if err := db_helpers.NewSeeder(
		db,
		context.TODO(),
	).SeedChatsN(1); err != nil {
		t.Fatalf("failed to seed latest chat: %s\n", err)
	}
	coreInstance.GetConfig().DatabaseDSN = "deez"
	cmd := chat.CreateChatCommand(coreInstance, test_helpers.MockCompletion)
	cmd.Flags().Set("content", "testing stuff")
	cmd.Flags().Set("latest", "true")

	relayDataChan := make(chan string)
	requestIDChan := make(chan string)
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		id := r.Header.Get("ID")
		dets := fmt.Sprintf("%s | %s | %s", r.Method, r.URL.String(), body)
		go func() {
			relayDataChan <- dets
			requestIDChan <- id
		}()
	}))
	server := httptest.NewServer(mux)
	defer server.Close()

	db_helpers.CreateActiveSession(
		db,
		context.Background(),
		&models.ActiveSession{
			DatabaseDNS: coreInstance.GetConfig().DatabaseDSN,
			Address:     server.URL,
		},
	)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute cmd: %s", err)
	}

	relayData := []string{}
	idsData := []string{}
	for i := 0; i < 4; i++ {
		select {
		case data := <-relayDataChan:
			relayData = append(relayData, data)
		case data := <-requestIDChan:
			idsData = append(idsData, data)
		}
	}

	dbMessages, err := db_helpers.GetMessagesByChatID(db, context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to retrieve db chat messages: %s", err)
	}

	dbMessages[0].ID = 0                     // Optimistic
	dbMessages[0].TimestampsToNilForTest__() // Optimistic
	serverChatCreatedMessage := messages.ServerMessageCreated{
		ID:   idsData[0],
		Type: "message-created",
		Payload: messages.ServerMessageCreatedPayload{
			ChatID:  1,
			Message: dbMessages[0],
		},
	}

	chatCreatedData, err := messages.Encode(serverChatCreatedMessage)
	if err != nil {
		t.Fatalf("failed to encode expected chat message for test: %s", err)
	}

	expectedChatCreatedData := fmt.Sprintf("POST | /api/v1/relay | %s", chatCreatedData)
	if relayData[0] != expectedChatCreatedData {
		t.Fatalf(
			"failed to get chat create message event\nexpected: %v\nactual:   %v\n",
			expectedChatCreatedData,
			relayData[0],
		)
	}

	serverMessageCreated := messages.ServerMessageCreated{
		ID:   idsData[1],
		Type: "message-created",
		Payload: messages.ServerMessageCreatedPayload{
			ChatID:  1,
			Message: dbMessages[1],
		},
	}

	messageData, err := messages.Encode(serverMessageCreated)
	if err != nil {
		t.Fatalf("failed to encode expected chat message for test: %s", err)
	}

	expectedMessageData := fmt.Sprintf("POST | /api/v1/relay | %s", messageData)
	if relayData[1] != expectedMessageData {
		t.Fatalf(
			"failed to get chat create message event\nexpected: %v\nactual:   %v\n",
			expectedMessageData,
			relayData[1],
		)
	}

	out := coreInstance.GetConfig().Stdoout.(*strings.Builder)
	if out.String() != dbMessages[len(dbMessages)-1].Content+"\n" {
		t.Fatalf(
			"failed to write result to stdout\nexpected: %q\nactual:   %q\n",
			"> mocked: testing stuff\n",
			out.String(),
		)
	}
}
