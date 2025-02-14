package template_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k10wl/hermes/cmd/template"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
)

func TestUpsertTemplate(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	ctx := context.Background()
	cmd := template.CreateTemplateCommand(coreInstance)
	content := `--{{define "custom"}}[--{{.}}]--{{end}}`
	cmd.SetArgs([]string{"upsert", "--content", content})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute cmd: %s\n", err)
	}
	tmp, err := db_helpers.GetTemplateByName(db, ctx, "custom")
	if err != nil {
		t.Fatalf("failed to retrieve deleted template - %s\n", err)
	}
	if tmp.Name != "custom" {
		t.Fatalf("failed to infer template name\nexpected: %q\actual:   %q\n", "custom", tmp.Name)
	}
	if tmp.Content != content {
		t.Fatalf("failed to store exact template content\nexpected: %q\actual:   %q\n", content, tmp.Content)
	}
}

func TestUpsertRelay(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	cmd := template.CreateTemplateCommand(coreInstance)
	cmd.SetArgs([]string{"upsert", "--content", `--{{define "custom"}}[--{{.}}]--{{end}}`})
	ctx := context.Background()

	relayDataChan := make(chan string)
	requestIDChan := make(chan string)
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		id := r.Header.Get("ID")
		dets := fmt.Sprintf("%s | %s | %s", r.Method, r.URL.String(), body)
		go func() {
			relayDataChan <- dets
			requestIDChan <- id
		}()
	}))
	server := httptest.NewServer(mux)
	defer server.Close()
	db_helpers.CreateActiveSession(
		db,
		context.Background(),
		&models.ActiveSession{
			DatabaseDNS: coreInstance.GetConfig().DatabaseDSN,
			Address:     server.URL,
		},
	)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute cmd: %s\n", err)
	}

	relayData := ""
	idData := ""
	for i := 0; i < 2; i++ {
		select {
		case data := <-relayDataChan:
			relayData = data
		case data := <-requestIDChan:
			idData = data
		}
	}

	tmp, err := db_helpers.GetTemplateByName(db, ctx, "custom")
	if err != nil {
		t.Fatalf("failed to retrieve 'custom' template from database - %s\n", err)
	}

	byte, err := messages.Encode(messages.NewServerTemplateCreated(idData, tmp))
	if err != nil {
		t.Fatalf("failed to encode retrieved template to bytes - %s\n", err)
	}

	expected := fmt.Sprintf(`POST | /api/v1/relay | %s`, byte)
	if relayData != expected {
		t.Fatalf("expected relay message does not match with actual\nexpected: %q\nactual:   %q\n", expected, relayData)
	}
}
