package messages

import (
	"context"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
)

type ClientReadTemplatesPayload struct {
	StartBeforeID int64  `json:"start_before_id" validate:"required"`
	Limit         int64  `json:"limit"           validate:"required"`
	Name          string `json:"name"`
}

type ClientReadTemplates struct {
	ID      string                     `json:"id,required"      validate:"required,uuid4"`
	Type    string                     `json:"type,required"    validate:"required"`
	Payload ClientReadTemplatesPayload `json:"payload,required" validate:"required"`
}

func (message *ClientReadTemplates) Process(
	comms CommunicationChannel,
	c *core.Core,
	_ ai_clients.CompletionFn,
) error {
	cmd := core.NewGetTemplatesQuery(
		c,
		message.Payload.StartBeforeID,
		message.Payload.Limit,
		message.Payload.Name,
	)
	if err := cmd.Execute(context.Background()); err != nil {
		return BroadcastServerEmittedMessage(
			comms.Single(),
			NewServerError(message.ID, err.Error()),
		)
	}
	return BroadcastServerEmittedMessage(
		comms.Single(),
		NewServerReadTemplates(message.ID, cmd.Result),
	)
}

func (message *ClientReadTemplates) GetID() string { return message.ID }

type ClientReadTemplatePayload struct {
	ID int64 `json:"id" validate:"required"`
}

type ClientReadTemplate struct {
	ID      string                    `json:"id,required"      validate:"required,uuid4"`
	Type    string                    `json:"type,required"    validate:"required"`
	Payload ClientReadTemplatePayload `json:"payload,required" validate:"required"`
}

func (message *ClientReadTemplate) Process(
	comms CommunicationChannel,
	c *core.Core,
	_ ai_clients.CompletionFn,
) error {
	cmd := core.NewGetTemplateByIDQuery(
		c,
		message.Payload.ID,
	)
	if err := cmd.Execute(context.Background()); err != nil {
		return BroadcastServerEmittedMessage(
			comms.Single(),
			NewServerError(message.ID, err.Error()),
		)
	}
	return BroadcastServerEmittedMessage(
		comms.Single(),
		NewServerReadTemplate(message.ID, cmd.Result),
	)
}

func (message *ClientReadTemplate) GetID() string { return message.ID }

type ClientEditTemplatePayload struct {
	Name    string `json:"name,required"    validate:"required"`
	Content string `json:"content,required" validate:"required"`
	Clone   bool   `json:"clone"`
}

type ClientEditTemplate struct {
	ID      string                    `json:"id,required"      validate:"required,uuid4"`
	Type    string                    `json:"type,required"`
	Payload ClientEditTemplatePayload `json:"payload,required"`
}

func (message *ClientEditTemplate) Process(
	comms CommunicationChannel,
	c *core.Core,
	completionFn ai_clients.CompletionFn,
) error {
	cmd := core.NewEditTemplateByName(
		c,
		message.Payload.Name,
		message.Payload.Content,
		message.Payload.Clone,
	)
	if err := cmd.Execute(context.TODO()); err != nil {
		return BroadcastServerEmittedMessage(comms.Single(), NewServerError(
			message.ID,
			err.Error(),
		))
	}
	if message.Payload.Clone {
		return BroadcastServerEmittedMessage(
			comms.All(),
			NewServerTemplateCreated(message.ID, cmd.Result),
		)
	}
	return BroadcastServerEmittedMessage(
		comms.All(),
		NewServerTemplateChanged(message.ID, cmd.Result),
	)
}

func (message *ClientEditTemplate) GetID() string {
	return message.ID
}

type ClientDeleteTemplatePayload struct {
	Name string `json:"name,required"    validate:"required"`
}

type ClientDeleteTemplate struct {
	ID      string                      `json:"id,required"      validate:"required,uuid4"`
	Type    string                      `json:"type,required"`
	Payload ClientDeleteTemplatePayload `json:"payload,required"`
}

func (message *ClientDeleteTemplate) Process(
	comms CommunicationChannel,
	c *core.Core,
	completionFn ai_clients.CompletionFn,
) error {
	cmd := core.NewDeleteTemplateByName(
		c,
		message.Payload.Name,
	)
	if err := cmd.Execute(context.TODO()); err != nil {
		return BroadcastServerEmittedMessage(comms.Single(), NewServerError(
			message.ID,
			err.Error(),
		))
	}
	return BroadcastServerEmittedMessage(
		comms.All(),
		NewServerTemplateDeleted(message.ID, message.Payload.Name),
	)
}

func (message *ClientDeleteTemplate) GetID() string {
	return message.ID
}
