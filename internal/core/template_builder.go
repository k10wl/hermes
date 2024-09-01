package core

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/k10wl/hermes/internal/models"
)

type templateBuilder struct {
	core             *Core
	requiredTemplate string
	storedTemplates  map[string]*models.Template
}

func newTemplateBuilder(core *Core) *templateBuilder {
	return &templateBuilder{
		core:             core,
		requiredTemplate: "",
		storedTemplates:  map[string]*models.Template{"": {}},
	}
}

func (tb *templateBuilder) mustProcessTemplate(template string) {
	tb.requiredTemplate = template
}

func (tb *templateBuilder) process(
	ctx context.Context,
	input string,
) {
	inputTemplates, err := extractTemplates(input)
	if err != nil {
		return
	}
	if _, ok := tb.storedTemplates[tb.requiredTemplate]; !ok {
		inputTemplates = append(inputTemplates, tb.requiredTemplate)
	}
	inputTemplates = tb.removeStored(inputTemplates)
	if len(inputTemplates) == 0 {
		return
	}
	query := NewGetTemplatesByNamesQuery(tb.core, inputTemplates)
	if err := query.Execute(ctx); err != nil {
		return
	}
	for _, template := range query.Result {
		tb.storedTemplates[template.Name] = template
		tb.process(ctx, template.Content)
	}
}

func (tb templateBuilder) string() (string, error) {
	if !tb.hasRequired() {
		return "", fmt.Errorf("does not contain required template")
	}
	buf := &strings.Builder{}
	for _, template := range tb.storedTemplates {
		buf.WriteString(template.Content)
	}
	return buf.String(), nil
}

func (tb templateBuilder) hasRequired() bool {
	if tb.requiredTemplate == "" {
		return true
	}
	_, ok := tb.storedTemplates[tb.requiredTemplate]
	return ok
}

func (tb templateBuilder) removeStored(templates []string) []string {
	return slices.DeleteFunc(templates, func(name string) bool {
		_, ok := tb.storedTemplates[name]
		return ok
	})

}
