package util

import (
	"bytes"
	"math/rand"
	"strings"
	gt "text/template"
	"unicode"

	"github.com/pkg/errors"
)

var templateFuncMap = gt.FuncMap{
	"sliceTakeStr":    tplSliceTakeStr,
	"sliceTakeStrRnd": tplSliceTakeStrRnd,
}

func ParseAndApplyTextTemplate(template string, variables any, compact bool) (string, error) {
	if len(template) == 0 || !strings.Contains(template, "{{") {
		// Shortcut: Template has no variables
		return template, nil
	}

	tplName := "Template: " + template[:20] + "..."
	tpl, err := gt.New(tplName).Funcs(templateFuncMap).Parse(template)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse template")
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, variables)
	if err != nil {
		return "", errors.Wrap(err, "Failed to execute template")
	}

	rendered := buffer.String()
	if compact {
		var result strings.Builder
		inWhitespace := false

		for _, r := range rendered {
			if unicode.IsSpace(r) {
				if !inWhitespace {
					result.WriteRune(r)
					inWhitespace = true
				}
			} else {
				result.WriteRune(r)
				inWhitespace = false
			}
		}

		rendered = result.String()
	}

	return rendered, nil
}

func tplSliceTakeStr(slice []string, limit int) []string {
	if len(slice) <= limit {
		return slice
	} else {
		return slice[:limit-1]
	}
}

func tplSliceTakeStrRnd(slice []string, limit int) []string {
	if len(slice) <= limit {
		return slice
	}

	shuffled := make([]string, len(slice))
	copy(shuffled, slice)

	for i := range shuffled {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:limit]
}
