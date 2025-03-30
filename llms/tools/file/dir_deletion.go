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

type DirectoryDeletionArguments struct {
	Path Path `json:"path"`
}

type DirectoryDeletionResult struct {
	Path string `json:"path"`
}

var DirectoryDeletionDefinition = tools.BuiltinDefinition{
	Description: "Delete a directory (including all files ans subdirectories) on the user's system.",
	Parameter: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "The path to the directory to delete. Use '~' as placeholder for the user's home directory.",
			},
		},
		"additionalProperties": false,
		"required":             []string{"path"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs DirectoryDeletionArguments
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

		err = os.RemoveAll(path)
		if err != nil {
			return nil, fmt.Errorf("error deleting directory: %w", err)
		}

		absolutePath, err := filepath.Abs(path)
		if err != nil {
			slog.Warn("Error getting absolute path!", "error", err)
			absolutePath = path
		}

		return json.Marshal(DirectoryDeletionResult{
			Path: absolutePath,
		})
	},
}
