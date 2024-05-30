package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/k10wl/hermes/internal/core"
	hermes_runtime "github.com/k10wl/hermes/internal/runtime"
)

//go:embed assets
var assetsEmbed embed.FS

//go:embed views
var viewsEmbed embed.FS

func Serve(core *core.Core, config *hermes_runtime.Config) error {
	server := NewServer(core)
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	openBrowser(fmt.Sprintf("http://%s", addr))
	httpServer := http.Server{
		Addr:    addr,
		Handler: server,
	}
	fmt.Printf("Starting server on %s\n", addr)
	return httpServer.ListenAndServe()
}

func NewServer(core *core.Core) http.Handler {
	mux := http.NewServeMux()
	t := NewTemplate()
	addRoutes(mux, core, t)
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

func openBrowser(url string) {
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
