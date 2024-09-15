package web

import (
	"encoding/json"
)

type Message struct {
	Type    string `json:"type,required"`
	Payload any    `json:"payload,omitempty"`
}

func newMessage(messageType string, payload any) Message {
	return Message{Type: messageType, Payload: payload}
}

func (message Message) encode() ([]byte, error) {
	return json.Marshal(message)
}
