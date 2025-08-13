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

// NegFloat32PtrToNil converts a negative float32 pointed to by a pointer to nil.
// This is useful when working with database fields that can be NULL, but Go float32s
// don't have a natural NULL state (they're always either a valid number or 0).
//
// Usage:
//
//	var f *float32 = new(float32)
//	*f = -1.0
//	NegFloat32PtrToNil(&f)  // Now f is nil instead of pointing to -1.0
//
// Parameters:
//
//	f - A pointer to a pointer to a float32 that may need conversion
func NegFloat32PtrToNil(f **float32) {
	if f != nil && *f != nil && **f < float32(0) {
		*f = nil
	}
}
