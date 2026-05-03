package util

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
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
		"sliceRandomN": tplSliceRandomN,
		"fmtEnum":      tplFmtEnum,
		"oneliner":     tplOneliner,
		"indent":       tplIndent,
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

func tplSliceRandomN(input any, limit int) (any, error) {
	v := reflect.ValueOf(input)

	if k := v.Kind(); k != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice, got %v", k)
	}

	length := v.Len()
	if limit > length {
		limit = length
	}

	indices := rand.Perm(length)[:limit]

	out := reflect.MakeSlice(v.Type(), limit, limit)
	for i, idx := range indices {
		out.Index(i).Set(v.Index(idx))
	}

	return out.Interface(), nil
}

func tplFmtEnum(v any) string {
	str := fmt.Sprintf("%v", v)
	words := strings.ReplaceAll(str, "_", " ")
	return strings.ToLower(words)
}

func tplOneliner(v any) string {
	str := fmt.Sprintf("%v", v)
	str = strings.ReplaceAll(str, "\n", " ")
	return strings.Join(strings.Fields(str), " ")
}

func tplIndent(indentSize int, v string) string {
	if len(v) == 0 {
		return v
	}

	lines := strings.Split(v, "\n")
	for i := range lines {
		lines[i] = strings.Repeat(" ", indentSize) + lines[i]
	}
	return strings.Join(lines, "\n")
}
