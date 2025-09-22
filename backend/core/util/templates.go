package util

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	gt "text/template"

	"github.com/pkg/errors"
)

var templateFuncMap = gt.FuncMap{
	"sliceTakeStr":    tplSliceTakeStr,
	"sliceTakeStrRnd": tplSliceTakeStrRnd,
	"fmtEnum":         tplFormatEnum,
}

func ContainsTemplateVars(template string) bool {
	return len(template) > 0 && strings.Contains(template, "{{")
}

func ParseAndApplyTextTemplate(template string, variables any) (string, error) {
	if !ContainsTemplateVars(template) {
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

	return buffer.String(), nil
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

func tplFormatEnum(v interface{}) string {
	return strings.ToLower(strings.ReplaceAll(fmt.Sprintf("%v", v), "_", " "))
}
