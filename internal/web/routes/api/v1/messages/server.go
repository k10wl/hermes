package messages

import (
	"encoding/json"
	"fmt"

	"github.com/k10wl/hermes/internal/models"
)

type ServerEmittedMessage interface {
	Encode() ([]byte, error)
}

func Broadcast(channel chan []byte, message ServerEmittedMessage) error {
	if message == nil {
		return fmt.Errorf(
			"failed to get message, expected interface, but got nil",
		)
	}
	data, err := message.Encode()
	if err != nil {
		return err
	}
	channel <- data
	return nil
}

type ServerReload struct {
	Type string `json:"type,required"`
}

func NewServerReload() *ServerReload {
	return &ServerReload{Type: "reload"}
}

func (message ServerReload) Encode() ([]byte, error) {
	return json.Marshal(message)
}

type ServerError struct {
	Type    string `json:"type,required"`
	Payload string `json:"payload,omitempty"`
}

func NewServerError(info string) *ServerError {
	return &ServerError{Type: "server-error", Payload: info}
}

func (message ServerError) Encode() ([]byte, error) {
	return json.Marshal(message)
}

type ServerReadChatPayload struct {
	ChatID   int64             `json:"chat_id,required"`
	Messages []*models.Message `json:"messages,required"`
}

type ServerReadChat struct {
	Type    string                `json:"type,required"`
	Payload ServerReadChatPayload `json:"payload,omitempty"`
}

func NewServerReadChat(chatID int64, messages []*models.Message) *ServerReadChat {
	return &ServerReadChat{Type: "read-chat", Payload: ServerReadChatPayload{
		ChatID:   chatID,
		Messages: messages,
	}}
}

func (message ServerReadChat) Encode() ([]byte, error) {
	return json.Marshal(message)
}

type ServerPong struct {
	Type string `json:"type,required"`
}

func NewServerPong() *ServerPong {
	return &ServerPong{Type: "pong"}
}

func (message ServerPong) Encode() ([]byte, error) {
	return json.Marshal(message)
}

type ServerMessageCreatedPayload struct {
	ChatID  int64           `json:"chat_id,required"`
	Message *models.Message `json:"message,required"`
}

type ServerMessageCreated struct {
	Type    string                      `json:"type,required"`
	Payload ServerMessageCreatedPayload `json:"payload,required"`
}

func NewServerMessageCreated(chatID int64, message *models.Message) *ServerMessageCreated {
	return &ServerMessageCreated{
		Type: "message-created",
		Payload: ServerMessageCreatedPayload{
			ChatID:  chatID,
			Message: message,
		},
	}
}

func (message ServerMessageCreated) Encode() ([]byte, error) {
	return json.Marshal(message)
}

type ServerChatCreatedPayload struct {
	Chat    *models.Chat    `json:"chat,required"`
	Message *models.Message `json:"message,required"`
}

type ServerChatCreated struct {
	Type    string                   `json:"type,required"`
	Payload ServerChatCreatedPayload `json:"payload,required"`
}

func NewServerChatCreated(chat *models.Chat, message *models.Message) *ServerChatCreated {
	return &ServerChatCreated{
		Type: "chat-created",
		Payload: ServerChatCreatedPayload{
			Chat:    chat,
			Message: message,
		},
	}
}

func (message ServerChatCreated) Encode() ([]byte, error) {
	return json.Marshal(message)
}
