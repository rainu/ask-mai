package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/mcp/server/tools"
	"io"
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

var ChangeTimesTool = mcp.NewTool("changeTimes",
	mcp.WithDescription("Changes the access and/or modification time of a file or directory on the user's system."),
	mcp.WithString("path",
		mcp.Required(),
		mcp.Description("The path to the file or directory to change the times for. Use '~' as placeholder for the user's home directory."),
	),
	mcp.WithString("access_time",
		mcp.Description("The new access time of the file or directory. Use 'now' to set the current time. Otherwise the time in RFC3339 format."),
	),
	mcp.WithString("modification_time",
		mcp.Description("The new modification time of the file or directory. Use 'now' to set the current time. Otherwise the time in RFC3339 format."),
	),
)

var ChangeTimesToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var pArgs ChangeTimesArguments

	r, w := io.Pipe()
	go func() {
		defer w.Close()

		json.NewEncoder(w).Encode(request.Params.Arguments)
	}()

	err := json.NewDecoder(r).Decode(&pArgs)
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

	return mcp.NewToolResultText(""), nil
}
