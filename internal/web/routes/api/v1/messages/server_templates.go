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

type ServerReadTemplatePayload struct {
	Template *models.Template `json:"template,required"`
}

type ServerReadTemplate struct {
	ID      string                    `json:"id,required"       validate:"required,uuid4"`
	Type    string                    `json:"type,required"`
	Payload ServerReadTemplatePayload `json:"payload,omitempty"`
}

func NewServerReadTemplate(
	id string,
	template *models.Template,
) *ServerReadTemplate {
	return &ServerReadTemplate{
		ID:   id,
		Type: "read-template",
		Payload: ServerReadTemplatePayload{
			Template: template,
		}}
}

func (message ServerReadTemplate) __serverMessageSignature() {}
