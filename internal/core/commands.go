package core

import (
	"context"
	"fmt"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
	"github.com/k10wl/hermes/internal/models"
)

type Command interface {
	Execute(context.Context) error
}

type CreateChatAndCompletionCommand struct {
	Core     *Core
	Role     string
	Message  string
	Template string
	Result   *models.Message
}

func NewCreateChatAndCompletionCommand(
	Core *Core,
	Role string,
	Message string,
	Template string,
) *CreateChatAndCompletionCommand {
	return &CreateChatAndCompletionCommand{
		Core:     Core,
		Role:     Role,
		Message:  Message,
		Template: Template,
	}
}

func (c *CreateChatAndCompletionCommand) Execute(ctx context.Context) error {
	if err := c.Core.assertAI(); err != nil {
		return err
	}
	input, err := c.Core.prepareMessage(ctx, c.Message, c.Template)
	if err != nil {
		return err
	}
	chat, _, err := c.Core.db.CreateChatAndMessage(
		ctx,
		c.Role,
		input,
	)
	if err != nil {
		return err
	}
	// TODO insert used value into the db and adjust queries to receive less messages
	res, _, err := c.Core.ai_client.ChatCompletion(
		[]ai_clients.Message{{Content: input, Role: UserRole}},
	)
	if err != nil {
		return err
	}
	message, err := c.Core.db.CreateMessage(
		ctx,
		chat.ID,
		res.Role,
		res.Content,
	)
	c.Result = message
	return err
}

type CreateCompletionCommand struct {
	Core     *Core
	Message  string
	Template string
	Role     string
	ChatID   int64
	Result   *models.Message
}

func NewCreateCompletionCommand(
	Core *Core,
	ChatID int64,
	Role string,
	Message string,
	Template string,
) *CreateCompletionCommand {
	return &CreateCompletionCommand{
		Core:     Core,
		ChatID:   ChatID,
		Message:  Message,
		Template: Template,
		Role:     Role,
	}
}

func (c *CreateCompletionCommand) Execute(ctx context.Context) error {
	if err := c.Core.assertAI(); err != nil {
		return err
	}
	input, err := c.Core.prepareMessage(ctx, c.Message, c.Template)
	if err != nil {
		return err
	}
	prev, err := c.Core.db.GetChatMessages(ctx, c.ChatID)
	if err != nil {
		return err
	}
	_, err = c.Core.db.CreateMessage(
		ctx,
		c.ChatID,
		c.Role,
		input,
	)
	if err != nil {
		return err
	}
	history := []ai_clients.Message{}
	for _, p := range prev {
		history = append(history, messageToAIMessage(p))
	}
	history = append(history, ai_clients.Message{Content: input, Role: UserRole})
	// TODO insert used value into the db and adjust queries to receive less messages
	res, _, err := c.Core.ai_client.ChatCompletion(history)
	if err != nil {
		return err
	}
	message, err := c.Core.db.CreateMessage(
		ctx,
		c.ChatID,
		res.Role,
		res.Content,
	)
	c.Result = message
	return err
}

type UpdateWebSettingsCommand struct {
	Core        *Core
	WebSettings models.WebSettings
}

func NewUpdateWebSettingsCommand(
	Core *Core,
	WebSettings models.WebSettings,
) *UpdateWebSettingsCommand {
	return &UpdateWebSettingsCommand{
		Core:        Core,
		WebSettings: WebSettings,
	}
}

func (c *UpdateWebSettingsCommand) Execute(ctx context.Context) error {
	return c.Core.db.UpdateWebSettings(ctx, c.WebSettings.DarkMode)
}

type UpsertTemplateCommand struct {
	Core     *Core
	name     string
	template string
}

func NewUpsertTemplateCommand(core *Core, template string) *UpsertTemplateCommand {
	return &UpsertTemplateCommand{
		Core:     core,
		template: template,
	}
}

func (c UpsertTemplateCommand) Execute(ctx context.Context) error {
	name, err := extractTemplateDefinitionName(c.template)
	if err != nil {
		return err
	}
	_, err = c.Core.db.UpsertTemplate(ctx, name, c.template)
	return err
}

type DeleteTemplateByName struct {
	core *Core
	name string
}

func NewDeleteTemplateByName(core *Core, name string) *DeleteTemplateByName {
	return &DeleteTemplateByName{
		core: core,
		name: name,
	}
}

func (c DeleteTemplateByName) Execute(ctx context.Context) error {
	ok, err := c.core.db.DeleteTemplateByName(ctx, c.name)
	if !ok {
		return fmt.Errorf("did not remove any records")
	}
	return err
}
