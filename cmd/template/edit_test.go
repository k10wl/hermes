package template_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/k10wl/hermes/cmd/template"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
)

func TestEditTemplate(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	ctx := context.Background()
	cmd := template.CreateTemplateCommand(coreInstance)
	if err := db_helpers.CreateTemplate(db, ctx, &models.Template{
		Name:    "custom",
		Content: `--{{define "custom"}}[--{{.}}]--{{end}}`,
	}); err != nil {
		t.Fatalf("failed to create template for test - %s\n", err)
	}

	edit := `--{{define "custom"}}edited--{{end}}`
	cmd.SetArgs([]string{"edit", "--name", "custom", "--content", edit})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute cmd: %s\n", err)
	}

	tmp, err := db_helpers.GetTemplateByName(db, ctx, "custom")
	if err != nil {
		t.Fatalf("failed to retrieve deleted template - %s\n", err)
	}
	if tmp.Content != edit {
		t.Fatalf("failed to save edit\nexpected: %q\actual:   %q\n", edit, tmp.Content)
	}
}

func TestEditTemplateRename(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	ctx := context.Background()
	cmd := template.CreateTemplateCommand(coreInstance)
	if err := db_helpers.CreateTemplate(db, ctx, &models.Template{
		Name:    "custom",
		Content: `--{{define "custom"}}[--{{.}}]--{{end}}`,
	}); err != nil {
		t.Fatalf("failed to create template for test - %s\n", err)
	}

	edit := `--{{define "renamed"}}renamed--{{end}}`
	cmd.SetArgs([]string{"edit", "--name", "custom", "--content", edit})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute cmd: %s\n", err)
	}

	if _, err := db_helpers.GetTemplateByName(
		db,
		ctx,
		"custom",
	); !strings.Contains(err.Error(), "sql: no rows in result set") {
		t.Fatalf("old template is still accessible in database - %s\n", err)
	}

	tmp, err := db_helpers.GetTemplateByName(db, ctx, "renamed")
	if err != nil {
		t.Fatalf("failed to retrieve renamed template - %s\n", err)
	}
	if tmp.Content != edit {
		t.Fatalf("failed to save edit\nexpected: %q\actual:   %q\n", edit, tmp.Content)
	}
	if tmp.ID != 1 {
		t.Fatalf("stored edited template under unexpected id -%d\n", tmp.ID)
	}
}

func TestEditTemplateRenameWithClone(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	ctx := context.Background()
	cmd := template.CreateTemplateCommand(coreInstance)
	initial := `--{{define "custom"}}[--{{.}}]--{{end}}`
	if err := db_helpers.CreateTemplate(db, ctx, &models.Template{
		Name:    "custom",
		Content: initial,
	}); err != nil {
		t.Fatalf("failed to create template for test - %s\n", err)
	}

	edit := `--{{define "renamed"}}renamed--{{end}}`
	cmd.SetArgs([]string{"edit", "--name", "custom", "--content", edit, "--clone"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute cmd: %s\n", err)
	}

	customTmp, err := db_helpers.GetTemplateByName(
		db,
		ctx,
		"custom",
	)
	if err != nil {
		t.Fatalf("failed to retrieve custom template - %s\n", err)
	}
	if customTmp.Content != initial {
		t.Fatalf(
			"unexpected change in original template\nexpected: %s\nactual:   %s\n",
			initial,
			customTmp.Content,
		)
	}

	renamedTmp, err := db_helpers.GetTemplateByName(db, ctx, "renamed")
	if err != nil {
		t.Fatalf("failed to retrieve renamed template - %s\n", err)
	}
	if renamedTmp.Name != "renamed" {
		t.Fatalf(
			"failed to infer template name\nexpected: %s\nactual:   %s\n",
			"renamed",
			renamedTmp.Name,
		)
	}
	if renamedTmp.Content != edit {
		t.Fatalf(
			"failed to save edit\nexpected: %q\actual:   %q\n",
			edit,
			renamedTmp.Content,
		)
	}
}

func TestEditRelay(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	cmd := template.CreateTemplateCommand(coreInstance)
	ctx := context.Background()
	if err := db_helpers.CreateTemplate(db, ctx, &models.Template{
		Name:    "custom",
		Content: `--{{define "custom"}}[--{{.}}]--{{end}}`,
	}); err != nil {
		t.Fatalf("failed to create template for test - %s\n", err)
	}
	edit := `--{{define "edit"}}edit--{{end}}`
	cmd.SetArgs([]string{
		"edit",
		"--name", "custom",
		"--content", edit,
	})

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

	tmp, err := db_helpers.GetTemplateByName(db, ctx, "edit")
	if err != nil {
		t.Fatalf("failed to retrieve 'edit' template from database - %s\n", err)
	}

	byte, err := messages.Encode(messages.NewServerTemplateChanged(idData, tmp))
	if err != nil {
		t.Fatalf("failed to encode retrieved template to bytes - %s\n", err)
	}

	expected := fmt.Sprintf(`POST | /api/v1/relay | %s`, byte)
	if relayData != expected {
		t.Fatalf("expected relay message does not match with actual\nexpected: %q\nactual:   %q\n", expected, relayData)
	}
}

func TestEditRelayClone(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	cmd := template.CreateTemplateCommand(coreInstance)
	ctx := context.Background()
	if err := db_helpers.CreateTemplate(db, ctx, &models.Template{
		Name:    "custom",
		Content: `--{{define "custom"}}[--{{.}}]--{{end}}`,
	}); err != nil {
		t.Fatalf("failed to create template for test - %s\n", err)
	}
	edit := `--{{define "edit"}}edit--{{end}}`
	cmd.SetArgs([]string{
		"edit",
		"--name", "custom",
		"--content", edit,
		"--clone",
	})

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

	tmp, err := db_helpers.GetTemplateByName(db, ctx, "edit")
	if err != nil {
		t.Fatalf("failed to retrieve 'edit' template from database - %s\n", err)
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
