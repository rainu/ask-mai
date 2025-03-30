package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/llms/tools"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type StatsArguments struct {
	Path Path `json:"path"`
}

type StatsResult struct {
	Path        string    `json:"path"`
	IsDirectory bool      `json:"isDirectory"`
	IsRegular   bool      `json:"isRegular"`
	Permissions string    `json:"permissions"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"modTime"`
}

var StatsDefinition = tools.BuiltinDefinition{
	Description: "Get stats of a file or directory on the user's system.",
	Parameter: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "The path to the file or directory to get info for. Use '~' as placeholder for the user's home directory.",
			},
		},
		"additionalProperties": false,
		"required":             []string{"path"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs StatsArguments
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

		stats, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("error getting stats: %w", err)
		}

		absolutePath, err := filepath.Abs(path)
		if err != nil {
			slog.Warn("Error getting absolute path!", "error", err)
			absolutePath = path
		}

		return json.Marshal(StatsResult{
			Path:        absolutePath,
			IsDirectory: stats.IsDir(),
			IsRegular:   stats.Mode().IsRegular(),
			Permissions: stats.Mode().Perm().String(),
			Size:        stats.Size(),
			ModTime:     stats.ModTime(),
		})
	},
}
