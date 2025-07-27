package util

import (
	"os"
	"strings"
)

func SetStringFromEnvIfPresent(envVarName string, dest *string) {
	value, isSet := os.LookupEnv(envVarName)
	if isSet {
		*dest = value
	}
}

func SetSliceFromEnvIfPresent(envVarName string, dest *[]string) {
	value, isSet := os.LookupEnv(envVarName)
	if isSet {
		if strings.Contains(value, ",") {
			*dest = strings.Split(value, ",")
		} else {
			*dest = []string{value}
		}
	}
}
