package util

import "strings"

func EmptyStrToNil(str *string) *string {
	if str == nil || strings.TrimSpace(*str) == "" {
		return nil
	}

	return str
}
