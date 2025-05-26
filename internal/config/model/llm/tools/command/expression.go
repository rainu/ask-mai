package command

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/internal/expression"
)

type Expression string

type Variables struct {
	FunctionDefinition FunctionDefinition `json:"fd"`
	Arguments          string             `json:"args"`
}

func (c Expression) Validate() error {
	if len(c) == 0 {
		return nil
	}

	return expression.Precompile(string(c))
}

func (c Expression) CommandFn(fd FunctionDefinition) CommandFn {
	return func(ctx context.Context, args string) ([]byte, error) {
		result, err := expression.Run(ctx, string(c), Variables{
			FunctionDefinition: fd,
			Arguments:          args,
		}).AsByteArray()
		if err != nil {
			return nil, fmt.Errorf("error running expression: %w", err)
		}

		return result, nil
	}
}
