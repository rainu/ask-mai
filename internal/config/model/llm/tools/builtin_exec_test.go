package tools

import (
	"context"
	"encoding/json"
	"github.com/rainu/ask-mai/internal/expression"
	"github.com/rainu/ask-mai/internal/llms/tools/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuiltin_CommandExecution(t *testing.T) {
	toTest := NewCommandExecution()

	assert.True(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}

func TestBuiltin_CommandExecution_Approval(t *testing.T) {
	tests := []struct {
		expression string
		args       command.CommandExecutionArguments
		expected   bool
	}{
		{
			expression.VarNameContext + ".args.name == 'test'",
			command.CommandExecutionArguments{
				Name: "test",
			},
			true,
		},
		{
			expression.VarNameContext + `.args.name.endsWith('find') && ` + expression.VarNameContext + `.args.arguments.findIndex(a => a === "-exec") === -1`,
			command.CommandExecutionArguments{
				Name:      "find",
				Arguments: []string{"/"},
			},
			true,
		},
		{
			expression.VarNameContext + `.args.name.endsWith('find') && ` + expression.VarNameContext + `.args.arguments.findIndex(a => a === "-exec") === -1`,
			command.CommandExecutionArguments{
				Name:      "find",
				Arguments: []string{"/", "-exec", "rm", "{}", ";"},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			toTest := NewCommandExecution()
			toTest.Approval = tt.expression

			j, err := json.Marshal(tt.args)
			require.NoError(t, err)

			result := toTest.AsFunctionDefinition().NeedApproval(context.Background(), string(j))

			assert.Equal(t, tt.expected, result)
		})
	}
}
