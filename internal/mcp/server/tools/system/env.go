package system

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/mcp/server/tools"
	"os"
	"strings"
)

type EnvironmentArguments struct {
}

type EnvironmentResult struct {
	Environment map[string]string `json:"env"`
}

var EnvironmentDefinition = tools.BuiltinDefinition{
	Description: "Get all environment variables of the user's system.",
	Parameter: map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"additionalProperties": false,
		"required":             []string{},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs EnvironmentArguments
		err := json.Unmarshal([]byte(jsonArguments), &pArgs)
		if err != nil {
			return nil, fmt.Errorf("error parsing arguments: %w", err)
		}

		result := EnvironmentResult{
			Environment: map[string]string{},
		}

		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			result.Environment[pair[0]] = pair[1]
		}

		return json.Marshal(result)
	},
}

var EnvironmentTool = mcp.NewTool("getEnvironment",
	mcp.WithDescription("Get all environment variables of the user's system."),
)

var EnvironmentToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := EnvironmentResult{
		Environment: map[string]string{},
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		result.Environment[pair[0]] = pair[1]
	}

	raw, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(raw)), nil
}
