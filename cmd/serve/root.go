package serve

import (
	"github.com/k10wl/hermes/cmd/utils"
	"github.com/k10wl/hermes/internal/web"
	"github.com/spf13/cobra"
)

var ServeCommand = &cobra.Command{
	Use:   "serve",
	Short: "serve http client",
	Long:  "Serve as a HTTP web server.",
	Example: `$ hermes serve
$ hermes server --hostname 192.168.1.1 --port 8080`,
	PreRun: func(cmd *cobra.Command, args []string) {
		core := utils.GetCore(cmd)
		config := core.GetConfig()
		hostname, err := cmd.Flags().GetString("hostname")
		if err != nil {
			utils.LogError(config.Stderr, err)
		}
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			utils.LogError(config.Stderr, err)
		}
		config.Port = port
		config.Host = hostname
	},
	Run: func(cmd *cobra.Command, args []string) {
		core := utils.GetCore(cmd)
		config := core.GetConfig()
		err := web.Serve(core, config)
		if err != nil {
			utils.LogError(config.Stderr, err)
		}
	},
}

func init() {
	ServeCommand.Flags().StringP("hostname", "H", "127.0.0.1", "Specify the hostname")
	ServeCommand.Flags().StringP("port", "p", "8123", "Set the port")
}
