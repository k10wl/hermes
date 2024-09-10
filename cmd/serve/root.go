package serve

import (
	"github.com/k10wl/hermes/cmd/utils"
	"github.com/k10wl/hermes/internal/web"
	"github.com/spf13/cobra"
)

var (
	port     string
	hostname string
	open     bool
	latest   bool
)

var ServeCommand = &cobra.Command{
	Use:   "serve",
	Short: "Serve http client",
	Long:  "Serve as a HTTP web server.",
	Example: `$ hermes serve
$ hermes serve --hostname 192.168.1.1 --port 8080
$ hermes serve --open --latest`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		core := utils.GetCore(cmd)
		config := core.GetConfig()
		h, err := cmd.Flags().GetString("hostname")
		if err != nil {
			utils.LogError(config.Stderr, err)
			return err
		}
		p, err := cmd.Flags().GetString("port")
		if err != nil {
			utils.LogError(config.Stderr, err)
			return err
		}
		o, err := cmd.Flags().GetBool("open")
		if err != nil {
			utils.LogError(config.Stderr, err)
			return err
		}
		l, err := cmd.Flags().GetBool("latest")
		if err != nil {
			utils.LogError(config.Stderr, err)
			return err
		}
		if l && !o {
			utils.LogFail(config.Stderr, "cannot use --latest without --open")
			return err
		}
		port = p
		hostname = h
		open = o
		latest = l
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		core := utils.GetCore(cmd)
		config := core.GetConfig()
		if open {
			web.OpenBrowser(
				web.GetUrl(
					web.BuildAddr(hostname, port),
					core,
					config,
					latest,
				),
			)
		}
		err := web.Serve(core, config, hostname, port)
		if err != nil {
			utils.LogError(config.Stderr, err)
		}
	},
}

func init() {
	ServeCommand.Flags().StringP(
		"hostname",
		"H",
		"127.0.0.1",
		"specify the hostname",
	)
	ServeCommand.Flags().StringP(
		"port",
		"p",
		"8123",
		"set the port",
	)
	ServeCommand.Flags().BoolP("open", "o", false, "opens server in browser")
	ServeCommand.Flags().BoolP(
		"latest",
		"l",
		false,
		"will open latest chat if --open was provided",
	)
}
