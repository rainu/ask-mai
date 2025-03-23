package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Stats struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y StatsResult    `config:"-"`
	Z StatsArguments `config:"-"`
}

func (f Stats) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getStats",
		Description: "Get stats of a file or directory on the user's system.",
		CommandFn:   f.Command,
		Parameters: map[string]any{
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
		NeedsApproval: f.NeedsApproval,
	}
}

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

func (f Stats) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
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
}
