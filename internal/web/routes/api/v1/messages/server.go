package messages

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/settings"
	"github.com/k10wl/hermes/internal/validator"
)

type ServerEmittedMessage interface{ __serverMessageSignature() }

func BroadcastServerEmittedMessage(channel chan []byte, message ServerEmittedMessage) error {
	if message == nil {
		return fmt.Errorf(
			"failed to get message, expected interface, but got nil",
		)
	}
	data, err := Encode(message)
	if err != nil {
		return err
	}
	return BroadcastData(channel, data)
}

func BroadcastData(channel chan []byte, data []byte) error {
	channel <- data
	if config, err := settings.GetInstance(); err == nil {
		fmt.Fprintf(config.Stdoout, "    -<send>-> %s\n", data)
	}
	return nil
}

type ServerReload struct {
	ID   string `json:"id,required"   validate:"required,uuid4"`
	Type string `json:"type,required"`
}

func NewServerReload() *ServerReload {
	return &ServerReload{ID: uuid.New().String(), Type: "reload"}
}

func (message ServerReload) __serverMessageSignature() {}

type ServerError struct {
	ID      string `json:"id,required"       validate:"required,uuid4"`
	Type    string `json:"type,required"`
	Payload string `json:"payload,omitempty"`
}

func NewServerError(id string, info string) *ServerError {
	return &ServerError{ID: id, Type: "server-error", Payload: info}
}

func (message ServerError) __serverMessageSignature() {}

type ServerReadChatPayload struct {
	ChatID   int64             `json:"chat_id,required"`
	Messages []*models.Message `json:"messages,required"`
}

type ServerReadChat struct {
	ID      string                `json:"id,required"       validate:"required,uuid4"`
	Type    string                `json:"type,required"`
	Payload ServerReadChatPayload `json:"payload,omitempty"`
}

func NewServerReadChat(id string, chatID int64, messages []*models.Message) *ServerReadChat {
	return &ServerReadChat{
		ID:   id,
		Type: "read-chat",
		Payload: ServerReadChatPayload{
			ChatID:   chatID,
			Messages: messages,
		}}
}

func (message ServerReadChat) __serverMessageSignature() {}

type ServerPong struct {
	ID   string `json:"id,required"   validate:"required"`
	Type string `json:"type,required"`
}

func NewServerPong(id string) *ServerPong {
	return &ServerPong{ID: id, Type: "pong"}
}

func (message ServerPong) __serverMessageSignature() {}

type ServerMessageCreatedPayload struct {
	ChatID  int64           `json:"chat_id,required"`
	Message *models.Message `json:"message,required"`
}

type ServerMessageCreated struct {
	ID      string                      `json:"id,required"      validate:"required,uuid4"`
	Type    string                      `json:"type,required"`
	Payload ServerMessageCreatedPayload `json:"payload,required"`
}

func NewServerMessageCreated(
	id string,
	chatID int64,
	message *models.Message,
) *ServerMessageCreated {
	return &ServerMessageCreated{
		ID:   id,
		Type: "message-created",
		Payload: ServerMessageCreatedPayload{
			ChatID:  chatID,
			Message: message,
		},
	}
}

func (message ServerMessageCreated) __serverMessageSignature() {}

type ServerChatCreatedPayload struct {
	Chat    *models.Chat    `json:"chat,required"`
	Message *models.Message `json:"message,required"`
}

type ServerChatCreated struct {
	ID      string                   `json:"id,required"      validate:"required,uuid4"`
	Type    string                   `json:"type,required"`
	Payload ServerChatCreatedPayload `json:"payload,required"`
}

func NewServerChatCreated(
	id string,
	chat *models.Chat,
	message *models.Message,
) *ServerChatCreated {
	return &ServerChatCreated{
		ID:   id,
		Type: "chat-created",
		Payload: ServerChatCreatedPayload{
			Chat:    chat,
			Message: message,
		},
	}
}

func (message ServerChatCreated) __serverMessageSignature() {}

func Encode(serverMessage ServerEmittedMessage) ([]byte, error) {
	err := validator.Validate.Struct(serverMessage)
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(serverMessage)
}

type ServerTemplateChangedPayload struct {
	Template models.Template `json:"template" validate:"required"`
}

type ServerTemplateChanged struct {
	ID      string                       `json:"id,required"      validate:"required,uuid4"`
	Type    string                       `json:"type,required"`
	Payload ServerTemplateChangedPayload `json:"payload,required"`
}

func NewServerTemplateChanged(
	id string,
	template *models.Template,
) *ServerTemplateChanged {
	return &ServerTemplateChanged{
		ID:   id,
		Type: "template-changed",
		Payload: ServerTemplateChangedPayload{
			Template: *template,
		},
	}
}

func (message ServerTemplateChanged) __serverMessageSignature() {}

type ServerTemplateCreatedPayload struct {
	Template models.Template `json:"template" validate:"required"`
}

type ServerTemplateCreated struct {
	ID      string                       `json:"id,required"      validate:"required,uuid4"`
	Type    string                       `json:"type,required"`
	Payload ServerTemplateCreatedPayload `json:"payload,required"`
}

func NewServerTemplateCreated(
	id string,
	template *models.Template,
) *ServerTemplateCreated {
	return &ServerTemplateCreated{
		ID:   id,
		Type: "template-created",
		Payload: ServerTemplateCreatedPayload{
			Template: *template,
		},
	}
}

func (message ServerTemplateCreated) __serverMessageSignature() {}
