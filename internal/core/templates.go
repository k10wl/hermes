package core

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/k10wl/hermes/internal/models"
)

const rootTemplateName = "!!__root__!!"
const inPlaceTemplateName = "!!__execute_in_place__!!"
const leftDelim = "--{{"
const rightDelim = "}}"

func (c Core) prepareMessage(
	ctx context.Context,
	input string,
	templateName string,
) (string, error) {
	trimmed := trim(input)
	templateBuilder := newTemplateBuilder(&c)
	templateBuilder.mustProcessTemplate(templateName)
	templateBuilder.process(ctx, trimmed)
	templateString, err := templateBuilder.string()
	if err != nil {
		return trimmed, err
	}
	t := prepareTemplates(templateString)
	refinedInput := prepareInput(templateName, trimmed)
	buf := &strings.Builder{}
	err = executor(t, buf, refinedInput)
	return trim(buf.String()), err
}

const definitionError = "failed to get template name"

func getTemplateNames(content string) ([]string, error) {
	names := []string{}
	tmpl := template.New("")
	tmpl = tmpl.Delims(leftDelim, rightDelim)
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return names, err
	}
	for _, template := range tmpl.Templates() {
		name := template.Name()
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}

// Deprecated. Please replace with getTemplateNames
func extractTemplateDefinitionName(content string) (string, error) {
	defineRegexp := regexp.MustCompile(withDelims(`define "(?P<name>.*?)"`))
	i := defineRegexp.SubexpIndex("name")
	res := defineRegexp.FindStringSubmatch(content)
	if len(res) < i+1 {
		return "", fmt.Errorf("failed to get template name")
	}
	return res[i], nil
}

func extractTemplates(content string) ([]string, error) {
	templateRegexp := regexp.MustCompile(
		withDelims(`\s*?template\s+"(?P<name>.*?)".*?`),
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
	tmpl = tmpl.Delims(leftDelim, rightDelim)
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

func withInPlaceBlock(str string) string {
	return fmt.Sprintf(
		`%s%s%s`,
		withDelims(fmt.Sprintf(`block "%s" .`, inPlaceTemplateName)),
		str,
		withDelims("end"),
	)
}

func detectTemplateUsage(inputTemplate string) bool {
	regex := regexp.MustCompile(withDelims(`\s*template\s+"(?P<name>[^"]+?)"(.*?)?\s*`))
	i := regex.SubexpIndex("name")
	matches := regex.FindAllStringSubmatch(inputTemplate, -1)
	str := []string{}
	for _, match := range matches {
		if len(match) < i || slices.Contains(str, match[i]) {
			continue
		}
		str = append(str, match[i])
	}
	return len(str) > 0
}

func withTemplateDefinition(name string, content string) string {
	return fmt.Sprintf(
		`%s%s%s%s`,
		withDelims(fmt.Sprintf("define %q", name)),
		withDelims("."),
		withDelims("end"),
		content,
	)
}

func prepareTemplates(
	templates string,
) *template.Template {
	tmpl := template.New(rootTemplateName)
	tmpl = tmpl.Delims(leftDelim, rightDelim)
	tmpl = template.Must(tmpl.Parse(withTemplateDefinition(rootTemplateName, templates)))
	return tmpl
}

func prepareInput(templateName string, input string) string {
	if templateName == "" {
		return input
	}
	return withDelims(fmt.Sprintf("template %q %q", templateName, input))
}

func executor(t *template.Template, writer io.Writer, str string) error {
	if detectTemplateUsage(str) {
		updated, _ := t.Parse(withInPlaceBlock(str))
		buf := &strings.Builder{}
		err := updated.ExecuteTemplate(buf, inPlaceTemplateName, nil)
		if err != nil {
			return err
		}
		return executor(updated, writer, buf.String())
	}
	return t.ExecuteTemplate(writer, rootTemplateName, str)
}

func withDelims(content string) string {
	return fmt.Sprintf("%s%s%s", leftDelim, content, rightDelim)
}

func trim(input string) string {
	return strings.TrimSpace(input)
}
