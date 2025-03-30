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

type FileTempCreationArguments struct {
	Content    string     `json:"content"`
	Suffix     string     `json:"suffix"`
	Permission Permission `json:"permission"`
}

type FileTempCreationResult struct {
	Path    string `json:"path"`
	Written int    `json:"written"`
}

var FileTempCreationDefinition = tools.BuiltinDefinition{
	Description: "Creates a new temporary file on the user's system.",
	Parameter: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"content": map[string]any{
				"type":        "string",
				"description": "The content of the file.",
			},
			"suffix": map[string]any{
				"type":        "string",
				"description": "The suffix of the file.",
			},
			"permission": map[string]any{
				"type":        "string",
				"description": "The permission of the file. Default is 0644.",
			},
		},
		"additionalProperties": false,
		"required":             []string{},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs FileTempCreationArguments
		err := json.Unmarshal([]byte(jsonArguments), &pArgs)
		if err != nil {
			return nil, fmt.Errorf("error parsing arguments: %w", err)
		}

		perm, err := pArgs.Permission.Get(os.FileMode(0644))
		if err != nil {
			return nil, err
		}

		file, err := os.CreateTemp("", "ask-mai.*"+pArgs.Suffix)
		if err != nil {
			return nil, fmt.Errorf("error creating file: %w", err)
		}
		defer file.Close()

		if err = os.Chmod(file.Name(), perm); err != nil {
			return nil, fmt.Errorf("error setting file permission: %w", err)
		}

		absolutePath, err := filepath.Abs(file.Name())
		if err != nil {
			slog.Warn("Error getting absolute path!", "error", err)
			absolutePath = file.Name()
		}

		s, err := file.WriteString(pArgs.Content)
		if err != nil {
			return nil, fmt.Errorf("error writing to file: %w", err)
		}

		return json.Marshal(FileTempCreationResult{
			Path:    absolutePath,
			Written: s,
		})
	},
}
