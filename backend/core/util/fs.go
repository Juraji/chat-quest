package util

import (
	"fmt"
	"io"
	"io/fs"
)

func ReadFileAsString(fs fs.FS, path string) (string, error) {
	file, err := fs.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to read template file '%s': %w", path, err)
	}
	defer file.Close()

	// Read the entire content of the file into a string
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read content from template file '%s': %w", path, err)
	}

	// Convert the byte slice to a string
	return string(content), nil
}
