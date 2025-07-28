package util

import "strings"

// EmptyStrPtrToNil converts an empty string pointed to by a pointer to nil.
// This is useful when working with database fields that can be NULL, but Go strings
// don't have a natural NULL state (they're always either a valid string or "").
//
// Usage:
//
//	var str *string = new(string)
//	*str = ""
//	EmptyStrPtrToNil(&str)  // Now str is nil instead of pointing to ""
//
// Parameters:
//
//	str - A pointer to a pointer to a string that may need conversion
func EmptyStrPtrToNil(str **string) {
	if str != nil && *str != nil && strings.TrimSpace(**str) == "" {
		*str = nil
	}
}
