package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Path string

func (p Path) Get() (string, error) {
	if strings.HasPrefix(string(p), "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting user's home directory: %w", err)
		}
		return filepath.Join(home, string(p)[1:]), nil
	}
	return string(p), nil
}
