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
	Long: `Receives template, parses, verifies and saves content into database.
If -c (--content) flag is not provided - default text editor will be opened.
Name of the template will derive from definition/block name.
Template must comply golang text template rules.
NOTE: delimiters differ from golang text template:
    - left delimiter  - '--{{';
    - right delimiter - '}}';
`,
	Example: `$ hermes template upsert
$ hermes template upsert -c "--{{define "template"}}(instruction)--{{end}}
$ hermes template upsert --content "--{{define "template"}}(instruction)--{{end}} `,
	Run: func(cmd *cobra.Command, args []string) {
		c := utils.GetCore(cmd)
		config := c.GetConfig()
		content, err := cmd.Flags().GetString("content")
		if err != nil {
			fmt.Fprintf(config.Stderr, "%v\n", err)
			return
		}
		if content == "" {
			content, err = utils.OpenInEditor(
				upsertTemplate,
				config.Stdin,
				config.Stdoout,
				config.Stderr,
			)
			if content == upsertTemplate {
				utils.LogFail(
					config.Stderr,
					"do not save example template, make some changes",
				)
				return
			}
			if err != nil {
				utils.LogError(config.Stderr, err)
				return
			}
		}
		if err := core.NewUpsertTemplateCommand(
			c,
			content,
		).Execute(context.Background()); err != nil {
			utils.LogError(config.Stderr, err)
			return
		}
		fmt.Fprintf(config.Stdoout, "Template upserted successfully\n")
	},
}

func init() {
	upsertCommand.Flags().StringP(
		"content",
		"c",
		"",
		"template content, if not provided - text editor will be opened",
	)
}
