package web

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/settings"
	v1 "github.com/k10wl/hermes/internal/web/routes/api/v1"
)

//go:embed assets
var assetsEmbed embed.FS

//go:embed views
var viewsEmbed embed.FS

func Serve(core *core.Core, config *settings.Config, hostname string, port string) error {
	server := NewServer(core)
	addr := BuildAddr(hostname, port)
	httpServer := http.Server{
		Addr:    addr,
		Handler: server,
	}
	fmt.Fprintf(config.Stdoout, "Starting server on %s\n", addr)
	echan := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			echan <- err
		}
		close(echan)
	}()
	select {
	case err := <-echan:
		return err
	case <-config.ShutdownContext.Done():
		fmt.Fprintln(config.Stdoout, "Shutdown signal received")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return httpServer.Shutdown(ctx)
}

func NewServer(core *core.Core) http.Handler {
	mux := http.NewServeMux()
	t := NewTemplate()
	hub := v1.NewHub()
	go hub.Run()
	addRoutes(mux, core, hub, t)
	return mux
}

func NewTemplate() *template.Template {
	tmpl := template.New("main")
	templateContent, err := viewsEmbed.ReadFile("views/home.html")
	if err != nil {
		panic(err)
	}
	tmpl, err = tmpl.Parse(string(templateContent))
	if err != nil {
		panic(err)
	}
	return tmpl
}

func OpenBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		// unsupported OS
	}
	if cmd == nil {
		fmt.Println("Cannot open browser automatically, unsupported OS")
	}
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

func GetUrl(addr string, c *core.Core, config *settings.Config, latest bool) string {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("http://%s", addr))
	if !latest {
		return str.String()
	}
	q := core.LatestChatQuery{
		Core: c,
	}
	err := q.Execute(context.Background())
	if err != nil {
		fmt.Println("Cannot get latest chat")
		return str.String()
	}
	str.WriteString(fmt.Sprintf("/chats/%d", q.Result.ID))
	return str.String()
}

func BuildAddr(hostname string, port string) string {
	return fmt.Sprintf("%s:%s", hostname, port)
}
