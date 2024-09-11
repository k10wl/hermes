package core

import (
	"context"
	"fmt"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/models"
)

type Command interface {
	Execute(context.Context) error
}

type CreateChatAndCompletionCommand struct {
	core       *Core
	role       string
	message    string
	template   string
	parameters *ai_clients.Parameters
	completion ai_clients.CompletionFn
	Result     *models.Message
}

func NewCreateChatAndCompletionCommand(
	core *Core,
	role string,
	message string,
	template string,
	parameters *ai_clients.Parameters,
	completion ai_clients.CompletionFn,
) *CreateChatAndCompletionCommand {
	return &CreateChatAndCompletionCommand{
		core:       core,
		role:       role,
		message:    message,
		template:   template,
		parameters: parameters,
		completion: completion,
	}
}

func (c *CreateChatAndCompletionCommand) Execute(ctx context.Context) error {
	input, err := c.core.prepareMessage(ctx, c.message, c.template)
	if err != nil {
		return err
	}
	chat, _, err := c.core.db.CreateChatAndMessage(
		ctx,
		c.role,
		input,
	)
	if err != nil {
		return err
	}
	// TODO insert used value into the db and adjust queries to receive less messages
	res, err := c.completion(
		[]*ai_clients.Message{{Content: input, Role: UserRole}},
		c.parameters,
		&c.core.config.Providers,
	)
	if err != nil {
		return err
	}
	message, err := c.core.db.CreateMessage(
		ctx,
		chat.ID,
		res.Role,
		res.Content,
	)
	c.Result = message
	return err
}

type CreateCompletionCommand struct {
	core       *Core
	message    string
	template   string
	role       string
	chatID     int64
	parameters *ai_clients.Parameters
	completion ai_clients.CompletionFn
	Result     *models.Message
}

func NewCreateCompletionCommand(
	core *Core,
	chatID int64,
	role string,
	message string,
	template string,
	parameters *ai_clients.Parameters,
	completion ai_clients.CompletionFn,
) *CreateCompletionCommand {
	return &CreateCompletionCommand{
		core:       core,
		chatID:     chatID,
		message:    message,
		template:   template,
		role:       role,
		parameters: parameters,
		completion: completion,
	}
}

func (c *CreateCompletionCommand) Execute(ctx context.Context) error {
	input, err := c.core.prepareMessage(ctx, c.message, c.template)
	if err != nil {
		return err
	}
	prev, err := c.core.db.GetChatMessages(ctx, c.chatID)
	if err != nil {
		return err
	}
	_, err = c.core.db.CreateMessage(
		ctx,
		c.chatID,
		c.role,
		input,
	)
	if err != nil {
		return err
	}
	history := []*ai_clients.Message{}
	for _, p := range prev {
		history = append(history, messageToAIMessage(p))
	}
	history = append(history, &ai_clients.Message{Content: input, Role: UserRole})
	// TODO insert used value into the db and adjust queries to receive less messages
	res, err := c.completion(history, c.parameters, &c.core.config.Providers)
	if err != nil {
		return err
	}
	message, err := c.core.db.CreateMessage(
		ctx,
		c.chatID,
		res.Role,
		res.Content,
	)
	c.Result = message
	return err
}

type UpdateWebSettingsCommand struct {
	core        *Core
	WebSettings models.WebSettings
}

func NewUpdateWebSettingsCommand(
	core *Core,
	WebSettings models.WebSettings,
) *UpdateWebSettingsCommand {
	return &UpdateWebSettingsCommand{
		core:        core,
		WebSettings: WebSettings,
	}
}

func (c *UpdateWebSettingsCommand) Execute(ctx context.Context) error {
	return c.core.db.UpdateWebSettings(ctx, c.WebSettings.DarkMode)
}

type UpsertTemplateCommand struct {
	core     *Core
	name     string
	template string
}

func NewUpsertTemplateCommand(core *Core, template string) *UpsertTemplateCommand {
	return &UpsertTemplateCommand{
		core:     core,
		template: template,
	}
}

func (c UpsertTemplateCommand) Execute(ctx context.Context) error {
	name, err := extractTemplateDefinitionName(c.template)
	if err != nil {
		return err
	}
	_, err = c.core.db.UpsertTemplate(ctx, name, c.template)
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
		return fmt.Errorf("Failed. Template %q not found.", c.name)
	}
	return err
}

type EditTemplateByName struct {
	core    *Core
	name    string
	content string
	clone   bool
}

func NewEditTemplateByName(
	core *Core,
	name string,
	content string,
	clone bool,
) *EditTemplateByName {
	return &EditTemplateByName{
		core:    core,
		name:    name,
		content: content,
		clone:   clone,
	}
}

func (c EditTemplateByName) Execute(ctx context.Context) error {
	names, err := getTemplateNames(c.content)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("content does not contain templates")
	}
	if len(names) != 1 {
		return fmt.Errorf("content contains multiple templates")
	}
	newName := names[0]
	if newName != c.name && c.clone {
		return c.handleClone(ctx, newName, c.content)
	}
	return c.handleEdit(ctx, newName)
}

func (c EditTemplateByName) handleClone(ctx context.Context, newName string, content string) error {
	templatesQuery := NewGetTemplatesByNamesQuery(c.core, []string{newName})
	if err := templatesQuery.Execute(ctx); err != nil {
		return err
	}
	if len(templatesQuery.Result) != 0 {
		return fmt.Errorf("template with given name already exists")
	}
	return NewUpsertTemplateCommand(c.core, content).Execute(ctx)
}

func (c EditTemplateByName) handleEdit(ctx context.Context, newName string) error {
	if newName != c.name {
		return c.renameAndDelete(ctx)
	}
	ok, err := c.core.db.EditTemplateByName(ctx, c.name, c.content)
	if !ok {
		return fmt.Errorf("did not update template, please make sure it exists")
	}
	return err
}

func (c EditTemplateByName) renameAndDelete(ctx context.Context) error {
	if err := NewUpsertTemplateCommand(
		c.core,
		c.content,
	).Execute(ctx); err != nil {
		return err
	}
	return NewDeleteTemplateByName(c.core, c.name).Execute(ctx)
}
