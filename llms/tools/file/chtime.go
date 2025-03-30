package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/llms/tools"
	"os"
	"strings"
	"time"
)

type Time string

func (t Time) Get() (time.Time, error) {
	if string(t) == "" {
		return time.Time{}, nil
	}
	if strings.ToLower(string(t)) == "now" {
		return time.Now(), nil
	}
	return time.Parse(time.RFC3339, string(t))
}

type ChangeTimesArguments struct {
	Path             Path `json:"path"`
	AccessTime       Time `json:"access_time"`
	ModificationTime Time `json:"modification_time"`
}

type ChangeTimesResult struct {
}

var ChangeTimesDefinition = tools.BuiltinDefinition{
	Description: "Changes the access and/or modification time of a file or directory on the user's system.",
	Parameter: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "The path to the file or directory to change the times for. Use '~' as placeholder for the user's home directory.",
			},
			"access_time": map[string]any{
				"type":        "string",
				"description": "The new access time of the file or directory. Use 'now' to set the current time. Otherwise the time in RFC3339 format.",
			},
			"modification_time": map[string]any{
				"type":        "string",
				"description": "The new modification time of the file or directory. Use 'now' to set the current time. Otherwise the time in RFC3339 format.",
			},
		},
		"additionalProperties": false,
		"required":             []string{"path"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs ChangeTimesArguments
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

		if string(pArgs.AccessTime) == "" && string(pArgs.ModificationTime) == "" {
			return nil, fmt.Errorf("missing parameter: at least one of 'access_time' or 'modification_time' must be set")
		}

		at, err := pArgs.AccessTime.Get()
		if err != nil {
			return nil, err
		}
		mt, err := pArgs.ModificationTime.Get()
		if err != nil {
			return nil, err
		}

		err = os.Chtimes(path, at, mt)
		if err != nil {
			return nil, fmt.Errorf("error changing times: %w", err)
		}

		return json.Marshal(ChangeTimesResult{})
	},
}
