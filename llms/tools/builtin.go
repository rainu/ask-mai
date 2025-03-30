package tools

import "context"

type BuiltinDefinition struct {
	Description string
	Parameter   map[string]any
	Function    func(ctx context.Context, jsonArguments string) ([]byte, error)
}
