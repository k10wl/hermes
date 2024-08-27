package core

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/k10wl/hermes/internal/models"
)

func (c Core) prepareMessage(
	ctx context.Context,
	input string,
	template string,
) (string, error) {
	templateBuilder := newTemplateBuilder(&c)
	templateBuilder.mustProcessTemplate(template)
	templateBuilder.process(ctx, input)
	templates, err := templateBuilder.string()
	if err != nil {
		return input, err
	}
	fmt.Printf("input: %v\n", input)
	fmt.Printf("templates: %v\n", templates)
	res, err := execute(templates, "", input)
	fmt.Printf("res: %v\n", res)
	if err != nil {
		return input, err
	}
	if template == "" {
		return res, nil
	}
	if innerTemplates, err := extractTemplates(input); err != nil || len(innerTemplates) == 0 {
		return res, nil
	}
	res, err = execute(templates, template, res)
	return res, err
}

const definitionError = "failed to get template name"

func extractTemplateDefinitionName(content string) (string, error) {
	defineRegexp := regexp.MustCompile(`{{define "(?P<name>.*?)"}}`)
	i := defineRegexp.SubexpIndex("name")
	res := defineRegexp.FindStringSubmatch(content)
	if len(res) < i+1 {
		return "", fmt.Errorf("failed to get template name")
	}
	return res[i], nil
}

func extractTemplates(content string) ([]string, error) {
	templateRegexp := regexp.MustCompile(
		`{{\s*?template\s+"(?P<name>.*?)".*?}}`,
	)
	i := templateRegexp.SubexpIndex("name")
	templateNames := []string{}
	var err error
	for _, match := range templateRegexp.FindAllStringSubmatch(content, -1) {
		if len(match) < i+1 {
			err = fmt.Errorf("failed to get template name")
			break
		}
		m := match[templateRegexp.SubexpIndex("name")]
		templateNames = append(templateNames, m)
	}
	return templateNames, err
}

func execute(content string, name string, input string) (string, error) {
	tmpl, err := template.New("").Parse(content)
	if err != nil {
		return "", err
	}
	executeName := name
	if name == "" {
		extractedName, err := extractTemplateDefinitionName(content)
		if err != nil && err.Error() == definitionError {
			return input, nil
		}
		if err != nil {
			return "", err
		}
		executeName = extractedName
	}
	buf := &strings.Builder{}
	err = tmpl.ExecuteTemplate(buf, executeName, input)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func concat(templates []*models.Template) string {
	buf := &strings.Builder{}
	for _, template := range templates {
		buf.WriteString(template.Content)
	}
	return buf.String()
}
