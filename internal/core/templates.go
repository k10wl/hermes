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
) string {
	fmt.Println("")
	defer fmt.Println("")
	usedTemplates, err := extractTemplates(input)
	if err != nil {
		return input
	}
	if template != "" {
		usedTemplates = append(usedTemplates, template)
	}
	if len(usedTemplates) == 0 {
		return input
	}
	query := NewGetTemplatesByNamesQuery(&c, usedTemplates)
	if err := query.Execute(ctx); err != nil {
		return input
	}
	content := concat(query.Result)
	fmt.Printf("template: %v\n", template)
	fmt.Printf("input: %v\n", input)
	fmt.Printf("content: %v\n", content)
	// must go into template definition
	// must be defined to combine template with input
	fmt.Printf("templatesContents: %v\n", content)
	res, err := execute(content, "", input)
	fmt.Printf("res1: %v\n", res)
	if err != nil {
		return input
	}
	if template == "" {
		return res
	}
	res, err = execute(content, template, res)
	fmt.Printf("res2: %v\n", res)
	return res
}

/*
{{define "wrapper"}}wrapper -- {{.}} -- wrapper{{end}}
{{define "hello"}}hello world{{end}}

{{template "wrapper" {{template "hello"}} }}
{{
*/

func extractTemplateDefinitionName(content string) (string, error) {
	defineRegexp := regexp.MustCompile(`{{define "(?P<name>.*)"}}`)
	i := defineRegexp.SubexpIndex("name")
	res := defineRegexp.FindStringSubmatch(content)
	if len(res) < i+1 {
		return "", fmt.Errorf("failed to get template name")
	}
	return res[i], nil
}

func extractTemplates(content string) ([]string, error) {
	templateRegexp := regexp.MustCompile(
		`{{(\s+)?template(\s+)"(?P<name>.*?)"(\s(\.([A-z]+)?)+)?(\s+)?}}`,
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

func execute(content string, name string, input any) (string, error) {
	tmpl, err := template.New("").Parse(content)
	if err != nil {
		return "", err
	}
	executeName := name
	if name == "" {
		extractedName, err := extractTemplateDefinitionName(content)
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
	for i, template := range templates {
		sufix := "\n"
		if len(templates) == i+1 {
			sufix = ""
		}
		buf.WriteString(fmt.Sprintf("%s%s", template.Content, sufix))
	}
	return buf.String()
}
