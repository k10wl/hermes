package core

import (
	"fmt"

	ai_clients "github.com/k10wl/hermes/internal/ai-clients"
)

type Command interface {
	Execute() error
}

type SendMessageCommand struct {
	Core    *Core
	Message string
	Result  string
}

func (c *SendMessageCommand) Execute() error {
	if c.Core == nil || c.Core.ai_client == nil {
		return fmt.Errorf("ai client not set")
	}
	res, err := c.Core.ai_client.ChatCompletion(
		[]ai_clients.Message{{Content: c.Message, Role: UserRole}},
	)
	if err != nil {
		return err
	}
	c.Result = res.Content
	return nil
}
