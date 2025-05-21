package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/tools"
	"log/slog"
	"os"
	"path/filepath"
)

type FileDeletionArguments struct {
	Path Path `json:"path"`
}

type FileDeletionResult struct {
	Path string `json:"path"`
}

var FileDeletionDefinition = tools.BuiltinDefinition{
	Description: "Delete a file on the user's system.",
	Parameter: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "The path to the file to delete. Use '~' as placeholder for the user's home directory.",
			},
		},
		"additionalProperties": false,
		"required":             []string{"path"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs FileDeletionArguments
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

		absolutePath, err := filepath.Abs(path)
		if err != nil {
			slog.Warn("Error getting absolute path!", "error", err)
			absolutePath = path
		}

		err = os.Remove(path)
		if err != nil {
			return nil, fmt.Errorf("error deleting file: %w", err)
		}

		return json.Marshal(FileDeletionResult{
			Path: absolutePath,
		})
	},
}
