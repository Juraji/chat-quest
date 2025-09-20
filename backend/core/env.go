package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Environment struct {
	DataDirectory    string
	ChatQuestUIRoot  string
	DebugEnabled     bool
	TrustedProxies   []string
	CorsAllowOrigins []string
	ApplicationHost  string
	ApplicationPort  string
	ApiBasePath      string
	DefaultFSPerm    os.FileMode
}

// MkDataDir creates directories in the application data directory and returns the full path.
// If one of the elements looks like a file name (contains a '.'), it stops creating directories.
// The function ensures all intermediate directories are created with proper permissions.
func (e Environment) MkDataDir(elements ...string) string {
	path := e.DataDirectory

	for _, element := range elements {
		path = filepath.Join(path, element)

		if strings.ContainsRune(element, '.') {
			// Found the filename, should be the end of the path.
			break
		}

		// Create the directory using the default perm with the executable bit.
		if err := os.Mkdir(path, e.DefaultFSPerm|0111); err != nil && !os.IsExist(err) {
			panic(errors.Wrapf(err, "Error creating directory: %s", path))
		}
	}

	return path
}

var currentEnvironment Environment

func InitEnvironment() {
	var err error

	currentEnvironment = Environment{
		DataDirectory:    "./data",
		ChatQuestUIRoot:  "./browser",
		DebugEnabled:     false,
		CorsAllowOrigins: []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		ApplicationHost:  "localhost",
		ApplicationPort:  "8080",
		ApiBasePath:      "/api",
		DefaultFSPerm:    0644,
	}

	setStringFromEnvIfPresent("CHAT_QUEST_DATA_DIR", &currentEnvironment.DataDirectory)
	if currentEnvironment.DataDirectory, err = filepath.Abs(currentEnvironment.DataDirectory); err != nil {
		panic(errors.Wrapf(err, "Failed to get absolute Path of CHAT_QUEST_DATA_DIR: %s", currentEnvironment.DataDirectory))
	}
	if err = os.MkdirAll(currentEnvironment.DataDirectory, os.ModePerm); err != nil {
		panic(errors.Wrapf(err, "Failed to create data directory at %v", currentEnvironment.DataDirectory))
	}

	setStringFromEnvIfPresent("CHAT_QUEST_UI_ROOT", &currentEnvironment.ChatQuestUIRoot)
	if currentEnvironment.ChatQuestUIRoot, err = filepath.Abs(currentEnvironment.ChatQuestUIRoot); err != nil {
		panic(errors.Wrapf(err, "Failed to get absolute Path of CHAT_QUEST_UI_ROOT: %s", currentEnvironment.ChatQuestUIRoot))
	}

	setSliceFromEnvIfPresent("CHAT_QUEST_TRUSTED_PROXIES", &currentEnvironment.TrustedProxies)
	setSliceFromEnvIfPresent("CHAT_QUEST_CORS_ALLOW_ORIGINS", &currentEnvironment.CorsAllowOrigins)
	setStringFromEnvIfPresent("CHAT_QUEST_APPLICATION_HOST", &currentEnvironment.ApplicationHost)
	setStringFromEnvIfPresent("CHAT_QUEST_APPLICATION_PORT", &currentEnvironment.ApplicationPort)
	setStringFromEnvIfPresent("CHAT_QUEST_API_BASE_PATH", &currentEnvironment.ApiBasePath)

	var debugModeVal string
	setStringFromEnvIfPresent("CHAT_QUEST_DEBUG", &debugModeVal)
	currentEnvironment.DebugEnabled = debugModeVal == "true"
}

func Env() Environment {
	return currentEnvironment
}

func setStringFromEnvIfPresent(envVarName string, dest *string) {
	value, isSet := os.LookupEnv(envVarName)
	if isSet {
		*dest = value
	}
}

func setSliceFromEnvIfPresent(envVarName string, dest *[]string) {
	value, isSet := os.LookupEnv(envVarName)
	if isSet {
		if strings.Contains(value, ",") {
			*dest = strings.Split(value, ",")
		} else {
			*dest = []string{value}
		}
	}
}
