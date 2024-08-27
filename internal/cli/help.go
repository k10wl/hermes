package cli

import (
	"fmt"
)

func GetHelpString(version string) string {
	return fmt.Sprintf(
		`Hermes - Host-based Extensible Response Management System - v%s

Usage:  hermes -content "Hello world!"
        cat logs.txt | hermes -content "show errors"

Hermes is a tool for communication and management of AI chats by 
accessing underlying API via terminal

Example:

        $ echo "Who are you?" | hermes
        I am a language model AI designed to assist with answering 
        questions and providing information to the best of my
        knowledge and abilities.`,
		version,
	)
}
