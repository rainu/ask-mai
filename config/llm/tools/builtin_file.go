package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

type Permission string

func (p Permission) Get(defaultPerm os.FileMode) (os.FileMode, error) {
	if string(p) != "" {
		pi, pe := strconv.ParseInt(string(p), 8, 32)
		if pe != nil {
			return defaultPerm, fmt.Errorf("error parsing permissions: %w", pe)
		}
		return os.FileMode(pi), nil
	}

	return defaultPerm, nil
}
