package utils

import (
	"fmt"
	"io"

	"github.com/k10wl/hermes/internal/models"
)

func ListTemplates(templates []*models.Template, w io.Writer) error {
	if len(templates) == 0 {
		_, err := fmt.Fprintf(
			w,
			"No templates matched search.\nUse -h to get info of how to add templates.\n",
		)
		return err
	}
	_, err := fmt.Fprintf(w, "List of templates:\n\n")
	for _, template := range templates {
		if e := writeRow(w, template.Name, template.Content); e != nil {
			err = e
			break
		}
	}
	return err
}

func writeRow(w io.Writer, name string, description string) error {
	_, err := fmt.Fprintf(
		w,
		"[Name]    %s\n[Content] %s\n--------------------\n",
		name,
		description,
	)
	return err
}
