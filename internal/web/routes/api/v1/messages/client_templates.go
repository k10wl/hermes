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
