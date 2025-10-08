package util

import (
	"strings"
)

func HasPrefixCaseInsensitive(s, prefix string) bool {
	return len(s) >= len(prefix) &&
		strings.EqualFold(s[:len(prefix)], prefix)
}

func HasSuffixCaseInsensitive(s, suffix string) bool {
	return len(s) >= len(suffix) &&
		strings.EqualFold(s[len(s)-len(suffix):], suffix)
}

func StrAsPointer(s string) *string {
	return &s
}

func StrPtrOrDefault(s *string, def string) string {
	if s != nil {
		return *s
	}
	return def
}
