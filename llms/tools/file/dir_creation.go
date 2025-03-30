package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/llms/tools"
	"log/slog"
	"os"
	"path/filepath"
)

type DirectoryCreationArguments struct {
	Path       Path       `json:"path"`
	Permission Permission `json:"permission"`
}

type DirectoryCreationResult struct {
	Path string `json:"path"`
}

var DirectoryCreationDefinition = tools.BuiltinDefinition{
	Description: "Creates a new directory (including all missing parent directories) on the user's system.",
	Parameter: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "The path to the directory to create. Use '~' as placeholder for the user's home directory.",
			},
			"permission": map[string]any{
				"type":        "string",
				"description": "The permission of the directory. Default is 0755.",
			},
		},
		"additionalProperties": false,
		"required":             []string{"path"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs DirectoryCreationArguments
		err := json.Unmarshal([]byte(jsonArguments), &pArgs)
		if err != nil {
			return nil, fmt.Errorf("error parsing arguments: %w", err)
		}

		if string(pArgs.Path) == "" {
			return nil, fmt.Errorf("missing parameter: 'path'")
		}
		path, err := pArgs.Path.Get()
		if err != nil {
			return nil, err
		}

		// Check if directory already exists
		dirInfo, dirErr := os.Stat(path)
		if dirErr == nil {
			if !dirInfo.IsDir() {
				return nil, fmt.Errorf("path exists but is a file: %s", path)
			}
			return nil, fmt.Errorf("directory already exists: %s", path)
		}

		perm, err := pArgs.Permission.Get(os.FileMode(0644))
		if err != nil {
			return nil, err
		}

		err = os.MkdirAll(path, perm)
		if err != nil {
			return nil, fmt.Errorf("error creating directory: %w", err)
		}

		absolutePath, err := filepath.Abs(path)
		if err != nil {
			slog.Warn("Error getting absolute path!", "error", err)
			absolutePath = path
		}

		return json.Marshal(DirectoryCreationResult{
			Path: absolutePath,
		})
	},
}
