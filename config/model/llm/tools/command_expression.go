package tools

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/expression"
)

type CommandExpression string

type CommandVariables struct {
	FunctionDefinition FunctionDefinition `json:"fd"`
	Arguments          string             `json:"args"`
}

func (c CommandExpression) Validate() error {
	if len(c) == 0 {
		return nil
	}

	return expression.Precompile(string(c))
}

func (c CommandExpression) CommandFn(fd FunctionDefinition) CommandFn {
	return func(ctx context.Context, args string) ([]byte, error) {
		result, err := expression.Run(ctx, string(c), CommandVariables{
			FunctionDefinition: fd,
			Arguments:          args,
		}).AsByteArray()
		if err != nil {
			return nil, fmt.Errorf("error running expression: %w", err)
		}

		return result, nil
	}
}
