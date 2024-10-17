package chat

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
	"github.com/spf13/cobra"
)

func notifyActiveSessions(c *core.Core) {
	config := c.GetConfig()
	db := c.GetDB()
	activeSession, err := db.GetActiveSessionByDatabaseDNS(config.DatabaseDSN)
	if err != nil {
		return
	}
	req, err := http.NewRequest(
		"GET",
		"http://"+activeSession.Address+"/api/v1/update",
		nil,
	)
	if err != nil {
		return
	}
	if _, err = http.DefaultClient.Do(req); err != nil {
		db.RemoveActiveSession(activeSession) // unreachable, remove from active sessions
		return
	}
}

func CreateChatCommand(c *core.Core, completion ai_clients.CompletionFn) *cobra.Command {
	stdin := ""
	aiParameters := ai_clients.Parameters{}

	chatCommand := &cobra.Command{
		Use:   "chat [flags] (will error upon empty content)",
		Short: "Send chat message for completion",
		Long: `Sends messages for AI completion. You can provide your message directly with the ` + "`--content`" + ` flag or pipe in text. Options include model selection, randomness adjustment, and template usage.
`,
		Example: `$ cat crash.log | hermes chat
$ hermes chat --content "hello world"

$ cat crash.log | hermes chat --content "what happened here?"
$ hermes chat --latest --content "how can I fix that crash I send you before?"

$ git diff --cached | hermes chat --template commit --model claude-haiku

$ hermes chat \
    --model gpt-4o \
    --max-tokens 10 \
    --temperature 0.2 \
    --content "is there a security risk in this message?" < risky_message.json`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			preloadParams(cmd, &aiParameters)
			config := c.GetConfig()
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				return nil
			}
			p, err := io.ReadAll(config.Stdin)
			if err != nil {
				return err
			}
			stdin = string(p)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := cmd.Flags().GetString("content")
			if err != nil {
				return err
			}
			template, err := cmd.Flags().GetString("template")
			if err != nil {
				return err
			}
			ok, err := cmd.Flags().GetBool("latest")
			if err != nil {
				return err
			}
			var complete func(
				c *core.Core,
				aiParameters *ai_clients.Parameters,
				content string,
				template string,
				completion ai_clients.CompletionFn,
			) error
			if ok {
				complete = completeInChat
			} else {
				complete = createChatAndComplete
			}
			if stdin != "" {
				content = fmt.Sprintf("%s\n\n%s", stdin, content)
			}
			if strings.Trim(content, " \n\t") == "" {
				return fmt.Errorf("input message was empty")
			}
			err = complete(c, &aiParameters, content, template, completion)
			if err != nil {
				return err
			}
			notifyActiveSessions(c)
			return nil
		},
	}

	chatCommand.Flags().SortFlags = false
	chatCommand.Flags().StringP(
		"content",
		"c",
		"",
		"completion content, can be combined with stdin",
	)
	chatCommand.Flags().StringP(
		"template",
		"t",
		"",
		"name of predefined template to be applied (see `hermes template --help)",
	)
	chatCommand.Flags().BoolP(
		"latest",
		"l",
		false,
		"continues conversation in latest chat",
	)
	chatCommand.Flags().StringP("model", "m", "gpt-4o-mini", "completion model")
	chatCommand.Flags().
		String(
			"temperature",
			"",
			"degree of randomness of AI answer (higher number - more chaotic)",
		)
	chatCommand.Flags().String(
		"max-tokens",
		"",
		"maximum number of tokens used in output",
	)

	return chatCommand
}

func completeInChat(
	c *core.Core,
	aiParameters *ai_clients.Parameters,
	content string,
	template string,
	completion ai_clients.CompletionFn,
) error {
	config := c.GetConfig()
	chatQuery := core.LatestChatQuery{Core: c}
	if err := chatQuery.Execute(config.ShutdownContext); err != nil {
		return err
	}
	cmd := core.NewCreateCompletionCommand(
		c,
		chatQuery.Result.ID,
		core.UserRole,
		content,
		template,
		aiParameters,
		completion,
	)
	err := cmd.Execute(config.ShutdownContext)
	if err != nil {
		return err
	}
	outputMessage(config.Stdoout, cmd.Result)
	return nil
}

func createChatAndComplete(
	c *core.Core,
	aiParameters *ai_clients.Parameters,
	content string,
	template string,
	completion ai_clients.CompletionFn,
) error {
	config := c.GetConfig()
	cmd := core.NewCreateChatAndCompletionCommand(
		c,
		core.UserRole,
		content,
		template,
		aiParameters,
		ai_clients.Complete,
	)
	err := cmd.Execute(c.GetConfig().ShutdownContext)
	if err != nil {
		return err
	}
	outputMessage(config.Stdoout, cmd.Result)
	return nil
}

func outputMessage(w io.Writer, message *models.Message) {
	fmt.Fprintf(w, "%s\n", message.Content)
}

func preloadParams(cmd *cobra.Command, params *ai_clients.Parameters) error {
	model, err := cmd.Flags().GetString("model")
	if err != nil {
		return err
	}
	maxTokens, err := cmd.Flags().GetString("max-tokens")
	if err != nil {
		return err
	}
	temperature, err := cmd.Flags().GetString("temperature")
	if err != nil {
		return err
	}
	params.Model = model
	params.MaxTokens = readMaxTokes(maxTokens)
	params.Temperature = readTemperature(temperature)
	return nil
}

func readTemperature(t string) *float64 {
	f, err := strconv.ParseFloat(t, 64)
	if err != nil {
		return nil
	}
	f64 := float64(f)
	return &f64
}

func readMaxTokes(t string) *int64 {
	i, err := strconv.Atoi(t)
	if err != nil {
		return nil
	}
	i64 := int64(i)
	return &i64
}
