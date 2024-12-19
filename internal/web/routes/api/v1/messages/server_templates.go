package messages

import "github.com/k10wl/hermes/internal/models"

type ServerReadTemplatesPayload struct {
	Templates []*models.Template `json:"templates,required"`
}

type ServerReadTemplates struct {
	ID      string                     `json:"id,required"       validate:"required,uuid4"`
	Type    string                     `json:"type,required"`
	Payload ServerReadTemplatesPayload `json:"payload,omitempty"`
}

func NewServerReadTemplates(
	id string,
	templates []*models.Template,
) *ServerReadTemplates {
	return &ServerReadTemplates{
		ID:   id,
		Type: "read-templates",
		Payload: ServerReadTemplatesPayload{
			Templates: templates,
		}}
}

func (message ServerReadTemplates) __serverMessageSignature() {}
