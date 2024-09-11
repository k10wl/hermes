package template

import (
	"context"
	"fmt"

	"github.com/k10wl/hermes/cmd/utils"
	"github.com/k10wl/hermes/internal/core"
	"github.com/spf13/cobra"
)

const upsertTemplate = `--{{define "example name"}}
<Prompt>
This block defines an example upsert template.
Prompt XML tag is not required, it helps to ` + "`cat`" + ` this text in VIM

Quick info about template capabilities:
>>> --{{.}} - Prints entire template input. If empty - prints <no value>

>>> --{{with .}}
      --{{.}} - Will print temlate input only if it is not empty
  --{{end}}

>>> --{{if .isEnabled}}
      --{{.jsonKey}} - prints out specific json key
      This block runs if the condition is true.
  --{{end}}

>>> Usefull examles from [docs](https://pkg.go.dev/text/template#hdr-Actions)
</Prompt>
--{{end}}
`

var upsertCommand = &cobra.Command{
	Use:   `upsert`,
	Short: "Update an existing template or create a new one if it does not exist.",
	Long: `Receives template, parses, verifies, and saves content into the database. If the ` + "`--content`" + ` (` + "`-c`" + `) flag is not provided, the default text editor will be opened. The name of the template will derive from the definition/block name. The template must comply with Golang text template rules. NOTE: Delimiters differ from Golang text template:
    - Left delimiter - '--{{';
    - Right delimiter - '}}';`,
	Example: `$ hermes template upsert
$ hermes template upsert -c "--{{define "template"}}(instruction)--{{end}}
$ hermes template upsert --content "--{{define "template"}}(instruction)--{{end}} `,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := utils.GetCore(cmd)
		config := c.GetConfig()
		content, err := cmd.Flags().GetString("content")
		if err != nil {
			return err
		}
		if content == "" {
			content, err = utils.OpenInEditor(
				upsertTemplate,
				config.Stdin,
				config.Stdoout,
				config.Stderr,
			)
			if content == upsertTemplate {
				return fmt.Errorf("do not save example template, make some changes\n")
			}
			if err != nil {
				return err
			}
		}
		if err := core.NewUpsertTemplateCommand(
			c,
			content,
		).Execute(context.Background()); err != nil {
			return err
		}
		fmt.Fprintf(config.Stdoout, "Template upserted successfully\n")
		return nil
	},
}

func init() {
	upsertCommand.Flags().StringP(
		"content",
		"c",
		"",
		"template content; if not provided, the text editor will be opened",
	)
}
