package util

import (
	"bytes"
	"strings"
	gt "text/template"
)

func NewTextTemplate(name string, template string) (*gt.Template, error) {
	return gt.New(name).Parse(template)
}

func WriteToString(template *gt.Template, data any) string {
	var buffer bytes.Buffer
	err := template.Execute(&buffer, data)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

func HasTemplateVars(template string) bool {
	return len(template) > 0 && strings.Contains(template, "{{")
}
