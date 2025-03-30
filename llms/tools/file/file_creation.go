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

type FileCreationArguments struct {
	Path       Path       `json:"path"`
	Content    string     `json:"content"`
	Permission Permission `json:"permission"`
}

type FileCreationResult struct {
	Path    string `json:"path"`
	Written int    `json:"written"`
}

var FileCreationDefinition = tools.BuiltinDefinition{
	Description: "Creates a new file on the user's system.",
	Parameter: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "The path to the file to create. Use '~' as placeholder for the user's home directory.",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "The content of the file.",
			},
			"permission": map[string]any{
				"type":        "string",
				"description": "The permission of the file. Default is 0644.",
			},
		},
		"additionalProperties": false,
		"required":             []string{"path"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs FileCreationArguments
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

		// Check if file already exists
		fileInfo, fileErr := os.Stat(path)
		if fileErr == nil {
			if fileInfo.IsDir() {
				return nil, fmt.Errorf("path exists but is a directory: %s", path)
			}
			return nil, fmt.Errorf("file already exists: %s", path)
		}

		flag := os.O_WRONLY | os.O_CREATE

		perm, err := pArgs.Permission.Get(os.FileMode(0644))
		if err != nil {
			return nil, err
		}

		file, err := os.OpenFile(path, flag, perm)
		if err != nil {
			return nil, fmt.Errorf("error creating file: %w", err)
		}
		defer file.Close()

		absolutePath, err := filepath.Abs(file.Name())
		if err != nil {
			slog.Warn("Error getting absolute path!", "error", err)
			absolutePath = file.Name()
		}

		s, err := file.WriteString(pArgs.Content)
		if err != nil {
			return nil, fmt.Errorf("error writing to file: %w", err)
		}

		return json.Marshal(FileCreationResult{
			Path:    absolutePath,
			Written: s,
		})
	},
}
