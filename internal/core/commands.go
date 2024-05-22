package core

import (
	"fmt"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
)

type Command interface {
	Execute()
}

type SendMessage struct {
	core    *Core
	Message string
	Result  string
}

func (c *SendMessage) Execute() error {
	if c.core == nil || c.core.ai_client == nil {
		return fmt.Errorf("ai client not set")
	}
	res, err := c.core.ai_client.ChatCompletion(
		[]ai_clients.Message{{Content: c.Message, Role: UserRole}},
	)
	if err != nil {
		return err
	}
	c.Result = res.Content
	return nil
}

func (core *Core) NewSendMessageCommand(message string) *SendMessage {
	return &SendMessage{core: core, Message: message}
}
