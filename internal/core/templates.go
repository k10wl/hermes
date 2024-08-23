package core

import (
	"fmt"
	"regexp"
)

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
