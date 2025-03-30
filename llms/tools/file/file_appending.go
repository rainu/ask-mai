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

type FileAppendingArguments struct {
	Path    Path   `json:"path"`
	Content string `json:"content"`
}

type FileAppendingResult struct {
	Path    string `json:"path"`
	Written int    `json:"written"`
}

var FileAppendingDefinition = tools.BuiltinDefinition{
	Description: "Append content to an existing file on the user's system.",
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
		},
		"additionalProperties": false,
		"required":             []string{"path"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs FileAppendingArguments
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
		if fileErr != nil {
			return nil, fmt.Errorf("file does not exists: %s", path)
		}
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("path is a directory: %s", path)
		}

		flag := os.O_WRONLY | os.O_APPEND

		file, err := os.OpenFile(path, flag, os.FileMode(0644))
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
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

		return json.Marshal(FileAppendingResult{
			Path:    absolutePath,
			Written: s,
		})
	},
}
