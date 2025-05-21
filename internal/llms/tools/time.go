package tools

import (
	"context"
	"time"
)

var SystemTimeDefinition = BuiltinDefinition{
	Description: "Get the current system time.",
	Parameter: map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		return []byte(time.Now().String()), nil
	},
}
