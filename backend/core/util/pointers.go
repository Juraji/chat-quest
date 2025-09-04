package util

import "strings"

func EmptyStrToNil(str *string) *string {
	if str == nil || strings.TrimSpace(*str) == "" {
		return nil
	}

	return str
}

func ZeroFloat32ToNil(f float32) *float32 {
	if f <= 0 {
		return nil
	} else {
		return &f
	}
}

func ZeroIntToNil(i int) *int {
	if i <= 0 {
		return nil
	} else {
		return &i
	}
}
