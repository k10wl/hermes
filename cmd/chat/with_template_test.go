package chat_test

import (
	"context"
	"testing"

	"github.com/k10wl/hermes/cmd/chat"
	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestShouldUseTemplate(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	ctx := context.Background()
	if err := db_helpers.CreateTemplate(db, ctx, &models.Template{
		Name:    "template",
		Content: `--{{define "template"}}[--{{.}}]--{{end}}`,
	}); err != nil {
		t.Fatalf("failed to create template, error: %s\n", err)
	}
	cmd := chat.CreateChatCommand(coreInstance, test_helpers.MockCompletion)
	cmd.Flags().Set("content", "content")
	cmd.Flags().Set("template", "template")
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to execute cmd: %s", err)
	}
	dbMessages, err := db_helpers.GetMessagesByChatID(db, ctx, 1)
	if err != nil {
		t.Fatalf("failed to retrieve messages, error: %s\n", err)
	}
	if dbMessages[0].Content != "[content]" {
		t.Fatalf(
			"failed to apply template, actual stored data: %s\n",
			dbMessages[0].Content,
		)
	}

}
