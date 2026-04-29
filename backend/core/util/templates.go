package util

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	gt "text/template"

	"github.com/pkg/errors"
)

func ContainsTemplateVars(template string) bool {
	return len(template) > 0 && strings.Contains(template, "{{")
}

func ParseAndApplyTextTemplate(name string, template string, variables any) (string, error) {
	if !ContainsTemplateVars(template) {
		// Shortcut: Template has no variables
		return template, nil
	}

	var templateFuncMap = gt.FuncMap{
		"sliceTakeStr":    tplSliceTakeStr,
		"sliceTakeStrRnd": tplSliceTakeStrRnd,
		"fmtEnum":         tplFormatEnum,
		"oneliner":        tplOneliner,
	}

	tpl, err := gt.New(name).Funcs(templateFuncMap).Parse(template)
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
	}

	return slice[:limit-1]
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
	str := fmt.Sprintf("%v", v)
	words := strings.ReplaceAll(str, "_", " ")
	return strings.ToLower(words)
}

func tplOneliner(v interface{}) string {
	str := fmt.Sprintf("%v", v)
	str = strings.ReplaceAll(str, "\n", " ")
	return strings.Join(strings.Fields(str), " ")
}
