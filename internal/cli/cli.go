package cli

import (
	"context"
	"fmt"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/runtime"
)

const HelpString = `Hermes - Host-based Extensible Response Management System

Usage:  hermes -m "Hello world!"
        cat logs.txt | hermes -m "show errors"

Hermes is a tool for communication and management of AI chats by 
accessing underlying API via terminal

Example:

        $ echo "Who are you?" | hermes
        I am a language model AI designed to assist with answering 
        questions and providing information to the best of my
        knowledge and abilities.`

func CLI(c *core.Core, config *runtime.Config) error {
	sendMessage := core.CreateChatAndCompletionCommand{
		Core:    c,
		Message: config.Prompt,
		Role:    "user",
	}
	if err := sendMessage.Execute(context.Background()); err != nil {
		return err
	}
	fmt.Fprintf(config.Stdoout, "%s\n", sendMessage.Result.Content)
	return nil
}
