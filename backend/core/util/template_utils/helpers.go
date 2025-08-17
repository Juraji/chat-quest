package template_utils

import (
	"bytes"
	"strings"
	gt "text/template"
)

func NewTemplate(name string, template string, funcMap gt.FuncMap) (*gt.Template, error) {
	return gt.New(name).Funcs(funcMap).Parse(template)
}

func NewTemplateWithLazy(name string, template string) (*gt.Template, error) {
	return NewTemplate(name, template, LazyTemplateFuncMap())
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
