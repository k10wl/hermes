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
)

func TestDeleteTemplate(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	ctx := context.Background()
	if err := db_helpers.CreateTemplate(db, ctx, &models.Template{
		Name:    "custom",
		Content: `--{{define "custom"}}[--{{.}}]--{{end}}`,
	}); err != nil {
		t.Fatalf("failed to create template, error: %s\n", err)
	}
	cmd := template.CreateTemplateCommand(coreInstance)
	cmd.SetArgs([]string{"delete", "--name", "custom"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute cmd: %s\n", err)
	}
	tmp, err := db_helpers.GetTemplateByName(db, ctx, "custom")
	if err != nil {
		t.Fatalf("failed to retrieve deleted template - %s\n", err)
	}
	if tmp.DeletedAt == nil {
		t.Fatalf("expected deleted template to be empty - %+v\n", tmp)
	}
}

func TestDeleteWithoutName(t *testing.T) {
	coreInstance, _ := test_helpers.CreateCore()
	sb := &strings.Builder{}
	cmd := template.CreateTemplateCommand(coreInstance)
	cmd.SetOut(sb)
	cmd.SetErr(sb)

	cmd.SetArgs([]string{"delete"})
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected cmd to error if name was not provided\n")
	}

	cmd.SetArgs([]string{"delete", "--name"})
	err = cmd.Execute()
	if err == nil {
		t.Fatalf("expected cmd to error if name was not provided\n")
	}
}

func TestDeleteNonExistingTemplate(t *testing.T) {
	coreInstance, _ := test_helpers.CreateCore()
	sb := &strings.Builder{}
	cmd := template.CreateTemplateCommand(coreInstance)
	cmd.SetOut(sb)
	cmd.SetErr(sb)
	cmd.SetArgs([]string{"delete", "--name", "custom"})
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected cmd to error if template does not exist\n")
	}
	if !strings.Contains(err.Error(), `Template "custom" not found`) {
		t.Fatalf("failed to notify reason of failure - %q\n", err.Error())
	}
}

func TestDeleteRelay(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	ctx := context.Background()
	if err := db_helpers.CreateTemplate(db, ctx, &models.Template{
		Name:    "custom",
		Content: `--{{define "custom"}}[--{{.}}]--{{end}}`,
	}); err != nil {
		t.Fatalf("failed to create template, error: %s\n", err)
	}
	cmd := template.CreateTemplateCommand(coreInstance)
	cmd.SetArgs([]string{"delete", "--name", "custom"})

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

	expected := fmt.Sprintf(`POST | /api/v1/relay | {"id":%q,"type":"template-deleted","payload":{"name":"custom"}}`, idData)
	if relayData != expected {
		t.Fatalf("expected relay message does not match with actual\nexpected: %q\nactual:   %q\n", expected, relayData)
	}
}
