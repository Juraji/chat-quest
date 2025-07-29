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

// NegFloat64PtrToNil converts a negative float64 pointed to by a pointer to nil.
// This is useful when working with database fields that can be NULL, but Go float64s
// don't have a natural NULL state (they're always either a valid number or 0).
//
// Usage:
//
//	var f *float64 = new(float64)
//	*f = -1.0
//	EmptyFloat64PtrToNil(&f)  // Now f is nil instead of pointing to -1.0
//
// Parameters:
//
//	f - A pointer to a pointer to a float64 that may need conversion
func NegFloat64PtrToNil(f **float64) {
	if f != nil && *f != nil && **f < float64(0) {
		*f = nil
	}
}
