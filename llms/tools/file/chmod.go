package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/llms/tools"
	"os"
)

type ChangeModeArguments struct {
	Path       Path       `json:"path"`
	Permission Permission `json:"permission"`
}

type ChangeModeResult struct {
}

var ChangeModeDefinition = tools.BuiltinDefinition{
	Description: "Changes the mode of file or directory on the user's system.",
	Parameter: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "The path to the file or directory to change the mode for. Use '~' as placeholder for the user's home directory.",
			},
			"permission": map[string]any{
				"type":        "string",
				"description": "The permission of the file or directory.",
			},
		},
		"additionalProperties": false,
		"required":             []string{"path", "permission"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs ChangeModeArguments
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

		if string(pArgs.Permission) == "" {
			return nil, fmt.Errorf("missing parameter: 'permission'")
		}
		perm, err := pArgs.Permission.Get(os.FileMode(0000))
		if err != nil {
			return nil, err
		}

		err = os.Chmod(path, perm)
		if err != nil {
			return nil, fmt.Errorf("error changing mode: %w", err)
		}

		return json.Marshal(ChangeModeResult{})
	},
}
