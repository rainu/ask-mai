package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type EnvironmentArguments struct {
}

type EnvironmentResult struct {
	Environment map[string]string `json:"env"`
}

var EnvironmentDefinition = BuiltinDefinition{
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
