package messages

import (
	"context"
	"fmt"
	"time"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/settings"
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
	Process(
		CommunicationChannel,
		*core.Core,
		ai_clients.CompletionFn,
	) error
	GetID() string
}

func ReadMessage(data []byte) (ClientEmittedMessage, error) {
	if config, err := settings.GetInstance(); err == nil {
		fmt.Fprintf(config.Stdoout, "   <-<read>-  %s\n", data)
	}
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
	case "request-read-templates":
		msg = &ClientReadTemplates{}
	case "request-read-template":
		msg = &ClientReadTemplate{}
	case "request-edit-template":
		msg = &ClientEditTemplate{}
	case "delete-template":
		msg = &ClientDeleteTemplate{}
	}
	if msg == nil {
		return nil, fmt.Errorf("received unknown message type\n")
	}
	err = decode(msg, data)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

type ClientRequestReadChatPayload struct {
}

type ClientRequestReadChat struct {
	ID      string `json:"id,required"       validate:"required,uuid4"`
	Type    string `json:"type,required"`
	Payload int64  `json:"payload,omitempty"`
}

func (message *ClientRequestReadChat) Process(
	comms CommunicationChannel,
	c *core.Core,
	_ ai_clients.CompletionFn,
) error {
	cmd := core.GetChatMessagesQuery{
		Core:   c,
		ChatID: message.Payload,
	}
	err := cmd.Execute(context.TODO())
	if err != nil {
		return BroadcastServerEmittedMessage(
			comms.Single(),
			NewServerError(
				message.ID,
				err.Error(),
			),
		)
	}
	return BroadcastServerEmittedMessage(
		comms.Single(),
		NewServerReadChat(message.ID, cmd.ChatID, cmd.Result),
	)
}

func (message *ClientRequestReadChat) GetID() string { return message.ID }

type ClientPing struct {
	ID   string `json:"id,required"   validate:"required,uuid4"`
	Type string `json:"type,required"`
}

func (message *ClientPing) Process(
	comms CommunicationChannel,
	_ *core.Core,
	_ ai_clients.CompletionFn,
) error {
	return BroadcastServerEmittedMessage(comms.Single(), NewServerPong(message.ID))
}

func (message *ClientPing) GetID() string { return message.ID }

type CreateCompletionMessagePayload struct {
	ChatID     int64                 `json:"chat_id"    validate:"required"`
	Content    string                `json:"content"    validate:"required"`
	Template   string                `json:"template"`
	Parameters ai_clients.Parameters `json:"parameters" validate:"required"`
}

type ClientCreateCompletion struct {
	ID      string                         `json:"id,required"      validate:"required,uuid4"`
	Type    string                         `json:"type,required"`
	Payload CreateCompletionMessagePayload `json:"payload,required"`
}

func (message *ClientCreateCompletion) GetID() string {
	return message.ID
}

func (message *ClientCreateCompletion) Process(
	comms CommunicationChannel,
	c *core.Core,
	completionFn ai_clients.CompletionFn,
) error {
	var fn func(
		comms CommunicationChannel,
		c *core.Core,
		completionFn ai_clients.CompletionFn,
	) error
	if message.Payload.ChatID == -1 {
		fn = message.processNewChat
	} else {
		fn = message.processExistingChat
	}
	return fn(comms, c, completionFn)
}

func (message *ClientCreateCompletion) processNewChat(
	comms CommunicationChannel,
	c *core.Core,
	completionFn ai_clients.CompletionFn,
) error {
	cmd := core.NewCreateChatWithMessageCommand(c, &models.Message{
		Role:    "user",
		Content: message.Payload.Content,
	}, "")
	if err := cmd.Execute(context.TODO()); err != nil {
		return err
	}
	if err := BroadcastServerEmittedMessage(
		comms.All(),
		NewServerChatCreated(message.ID, cmd.Result.Chat, cmd.Result.Message),
	); err != nil {
		return err
	}
	return message.createCompletion(
		comms,
		c,
		completionFn,
		cmd.Result.Chat.ID,
		false,
	)
}

func (message *ClientCreateCompletion) processExistingChat(
	comms CommunicationChannel,
	c *core.Core,
	completionFn ai_clients.CompletionFn,
) error {
	if err := BroadcastServerEmittedMessage(
		comms.All(),
		NewServerMessageCreated(
			message.ID,
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
		comms,
		c,
		completionFn,
		message.Payload.ChatID,
		true,
	)
}

func (message *ClientCreateCompletion) createCompletion(
	comms CommunicationChannel,
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
	cmd.ShouldPersistUserMessage(skipPersistingUserMessage)
	if err := cmd.Execute(context.TODO()); err != nil {
		return BroadcastServerEmittedMessage(comms.Single(), NewServerError(
			message.ID,
			err.Error(),
		))
	}
	return BroadcastServerEmittedMessage(
		comms.All(),
		NewServerMessageCreated(message.ID, chatID, cmd.Result),
	)
}
