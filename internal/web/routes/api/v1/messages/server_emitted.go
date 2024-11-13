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

type ReloadMessage struct {
	Type string `json:"type,required"`
}

func NewReloadMessage() *ReloadMessage {
	return &ReloadMessage{Type: "reload"}
}

func (message ReloadMessage) Encode() ([]byte, error) {
	return json.Marshal(message)
}

type ErrorMessage struct {
	Type    string `json:"type,required"`
	Payload string `json:"payload,omitempty"`
}

func NewErrorMessage(info string) *ErrorMessage {
	return &ErrorMessage{Type: "server-error", Payload: info}
}

func (message ErrorMessage) Encode() ([]byte, error) {
	return json.Marshal(message)
}

type ReadChatMessagePayload struct {
	ChatID   int64             `json:"chat_id,required"`
	Messages []*models.Message `json:"messages,required"`
}

type ReadChatMessage struct {
	Type    string                 `json:"type,required"`
	Payload ReadChatMessagePayload `json:"payload,omitempty"`
}

func NewReadChatMessage(chatID int64, messages []*models.Message) *ReadChatMessage {
	return &ReadChatMessage{Type: "read-chat", Payload: ReadChatMessagePayload{
		ChatID:   chatID,
		Messages: messages,
	}}
}

func (message ReadChatMessage) Encode() ([]byte, error) {
	return json.Marshal(message)
}

type PongMessage struct {
	Type string `json:"type,required"`
}

func NewPongMessage() *PongMessage {
	return &PongMessage{Type: "pong"}
}

func (message PongMessage) Encode() ([]byte, error) {
	return json.Marshal(message)
}
