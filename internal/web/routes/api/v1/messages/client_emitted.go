package messages

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/k10wl/hermes/internal/core"
)

const (
	None = iota
	Sender
	All
)

type ClientEmittedMessage interface {
	Decode(data []byte) error
	Process(core *core.Core) (message ServerEmittedMessage, receivers int, err error)
}

func typeDetector(data []byte) (string, error) {
	type typeDetector struct {
		Type string `json:"type,required"`
	}
	var t typeDetector
	err := json.Unmarshal(data, &t)
	return t.Type, err
}

func ReadMessage(data []byte) (ClientEmittedMessage, error) {
	messageType, err := typeDetector(data)
	if err != nil {
		return nil, err
	}
	switch messageType {
	case "ping":
		msg := &PingMessage{}
		return msg, msg.Decode(data)
	case "request-read-chat":
		msg := &RequestReadChatMessage{}
		return msg, msg.Decode(data)
	}
	return nil, fmt.Errorf("unhandled message type")
}

type RequestReadChatMessage struct {
	Type    string `json:"type,required"`
	Payload int64  `json:"payload,omitempty"`
}

func (message *RequestReadChatMessage) Decode(data []byte) error {
	return json.Unmarshal(data, message)
}

func (message *RequestReadChatMessage) Process(c *core.Core) (
	ServerEmittedMessage, int, error,
) {
	cmd := core.GetChatMessagesQuery{
		Core:   c,
		ChatID: message.Payload,
	}
	err := cmd.Execute(context.TODO())
	if err != nil {
		return nil, Sender, err
	}
	return NewReadChatMessage(cmd.ChatID, cmd.Result), Sender, nil
}

type PingMessage struct {
	Type string `json:"type,required"`
}

func (message *PingMessage) Decode(data []byte) error {
	return json.Unmarshal(data, message)
}

func (message *PingMessage) Process(_ *core.Core) (
	ServerEmittedMessage, int, error,
) {
	return NewPongMessage(), Sender, nil
}
