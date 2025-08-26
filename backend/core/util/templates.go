package util

import (
	"bytes"
	"github.com/pkg/errors"
	"strings"
	gt "text/template"
)

func ParseAndApplyTextTemplate(
	template string,
	variables any,
) (string, error) {
	if len(template) == 0 || !strings.Contains(template, "{{") {
		// Shortcut: Template has no variables
		return template, nil
	}

	tpl, err := gt.New("Template").Parse(template)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse template")
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, variables)
	if err != nil {
		return "", errors.Wrap(err, "Failed to execute template")
	}

	return buffer.String(), nil
}
