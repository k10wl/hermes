package serve

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/web"
	"github.com/spf13/cobra"
)

func CreateServeCommand(c *core.Core) *cobra.Command {
	var (
		port     string
		hostname string
		addr     string
		open     bool
		latest   bool
	)

	config := c.GetConfig()
	db := c.GetDB()
	activeSession := models.ActiveSession{}

	serveCommand := &cobra.Command{
		Use:   "serve",
		Short: "Serve http client",
		Long:  "Serve as a HTTP web server.",
		Example: `$ hermes serve
$ hermes serve --hostname 192.168.1.1 --port 8080
$ hermes serve --open --latest`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			h, err := cmd.Flags().GetString("hostname")
			if err != nil {
				return err
			}
			p, err := cmd.Flags().GetString("port")
			if err != nil {
				return err
			}
			o, err := cmd.Flags().GetBool("open")
			if err != nil {
				return err
			}
			l, err := cmd.Flags().GetBool("latest")
			if err != nil {
				return err
			}
			port = p
			hostname = h
			open = o
			latest = l
			addr = web.BuildAddr(h, p)
			activeSession.Address = addr
			activeSession.DatabaseDNS = config.DatabaseDSN
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("post run\n")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if latest && !open {
				return fmt.Errorf("cannot use --latest without --open\n")
			}
			if open {
				web.OpenBrowser(
					web.GetUrl(
						addr,
						c,
						config,
						latest,
					),
				)
			}
			err := db.CreateActiveSession(&activeSession)
			if err != nil {
				fmt.Fprintf(config.Stderr, "failed to store active session record - %s\n", err)
			}
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
			err = web.Serve(c, config, hostname, port)
			if err != nil {
				return err
			}
			<-quit
			go db.RemoveActiveSession(&activeSession)
			return nil
		},
	}

	serveCommand.Flags().SortFlags = false
	serveCommand.Flags().StringP(
		"hostname",
		"H",
		"127.0.0.1",
		"specify the hostname",
	)
	serveCommand.Flags().StringP(
		"port",
		"p",
		"8123",
		"set the port",
	)
	serveCommand.Flags().BoolP("open", "o", false, "opens server in browser")
	serveCommand.Flags().BoolP(
		"latest",
		"l",
		false,
		"will open latest chat if --open was provided",
	)

	return serveCommand
}
