package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
)

const (
	None = iota
	Sender
	All
)

type CommunicationChannel interface {
	Single() chan []byte
	All() chan []byte
}

type ClientEmittedMessage interface {
	Decode(data []byte) error
	Process(
		CommunicationChannel,
		*core.Core,
		ai_clients.CompletionFn,
	) error
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
		return nil, fmt.Errorf("failed to parse message type\n")
	}
	var msg ClientEmittedMessage
	switch messageType {
	case "ping":
		msg = &ClientPing{}
	case "create-completion":
		msg = &ClientCreateCompletion{}
	case "request-read-chat":
		msg = &ClientRequestReadChat{}
	}
	if msg == nil {
		return nil, fmt.Errorf("received unknown message type\n")
	}
	err = msg.Decode(data)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

type ClientRequestReadChat struct {
	Type    string `json:"type,required"`
	Payload int64  `json:"payload,omitempty"`
}

func (message *ClientRequestReadChat) Decode(data []byte) error {
	return json.Unmarshal(data, message)
}

func (message *ClientRequestReadChat) Process(
	coms CommunicationChannel,
	c *core.Core,
	_ ai_clients.CompletionFn,
) error {
	cmd := core.GetChatMessagesQuery{
		Core:   c,
		ChatID: message.Payload,
	}
	err := cmd.Execute(context.TODO())
	if err != nil {
		return Broadcast(
			coms.Single(),
			NewServerError(err.Error()),
		)
	}
	return Broadcast(
		coms.Single(),
		NewServerReadChat(cmd.ChatID, cmd.Result),
	)
}

type ClientPing struct {
	Type string `json:"type,required"`
}

func (message *ClientPing) Decode(data []byte) error {
	return json.Unmarshal(data, message)
}

func (message *ClientPing) Process(
	coms CommunicationChannel,
	_ *core.Core,
	_ ai_clients.CompletionFn,
) error {
	return Broadcast(coms.Single(), NewServerPong())
}

type CreateCompletionMessagePayload struct {
	ChatID     int64                 `json:"chat_id,required"`
	Content    string                `json:"content,required"`
	Template   string                `json:"template"`
	Parameters ai_clients.Parameters `json:"parameters,required"`
}

type ClientCreateCompletion struct {
	Type    string                         `json:"type,required"`
	Payload CreateCompletionMessagePayload `json:"payload,required"`
}

func (message *ClientCreateCompletion) Decode(data []byte) error {
	return json.Unmarshal(data, message)
}

func (message *ClientCreateCompletion) Process(
	coms CommunicationChannel,
	c *core.Core,
	completionFn ai_clients.CompletionFn,
) error {
	var fn func(
		coms CommunicationChannel,
		c *core.Core,
		completionFn ai_clients.CompletionFn,
	) error
	if message.Payload.ChatID == -1 {
		fn = message.processNewChat
	} else {
		fn = message.processExistingChat
	}
	return fn(coms, c, completionFn)
}

func (message *ClientCreateCompletion) processNewChat(
	coms CommunicationChannel,
	c *core.Core,
	completionFn ai_clients.CompletionFn,
) error {
	cmd := core.NewCreateChatWithMessageCommand(c, &models.Message{
		Role:    "user",
		Content: message.Payload.Content,
	})
	if err := cmd.Execute(context.TODO()); err != nil {
		return err
	}
	if err := Broadcast(
		coms.Single(),
		NewServerChatCreated(cmd.Result.Chat, cmd.Result.Message),
	); err != nil {
		return err
	}
	return message.createCompletion(
		coms,
		c,
		completionFn,
		cmd.Result.Chat.ID,
		false,
	)
}

func (message *ClientCreateCompletion) processExistingChat(
	coms CommunicationChannel,
	c *core.Core,
	completionFn ai_clients.CompletionFn,
) error {
	if err := Broadcast(
		coms.All(),
		NewServerMessageCreated(
			message.Payload.ChatID,
			&models.Message{
				ID:      time.Now().UnixMilli(),
				Content: message.Payload.Content,
				Role:    "user",
			},
		),
	); err != nil {
		return err
	}
	return message.createCompletion(
		coms,
		c,
		completionFn,
		message.Payload.ChatID,
		true,
	)
}

func (message *ClientCreateCompletion) createCompletion(
	coms CommunicationChannel,
	c *core.Core,
	completionFn ai_clients.CompletionFn,
	chatID int64,
	skipPersistingUserMessage bool,
) error {
	cmd := core.NewCreateCompletionCommand(
		c,
		chatID,
		"user",
		message.Payload.Content,
		message.Payload.Template,
		&message.Payload.Parameters,
		completionFn,
	)
	cmd.SkipPersistingUserMessage(skipPersistingUserMessage)
	if err := cmd.Execute(context.TODO()); err != nil {
		return Broadcast(coms.Single(), NewServerError(err.Error()))
	}
	return Broadcast(
		coms.All(),
		NewServerMessageCreated(chatID, cmd.Result),
	)
}
