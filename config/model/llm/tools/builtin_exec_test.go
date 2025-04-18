package tools

import (
	"github.com/rainu/ask-mai/expression"
	"github.com/rainu/ask-mai/llms/tools/command"
	"testing"
)

func TestCalcApprovalExpr(t *testing.T) {
	tests := []struct {
		expression  string
		args        command.CommandExecutionArguments
		expected    bool
		expectedErr error
	}{
		{
			expression.VarNameContext + ".name == 'test'",
			command.CommandExecutionArguments{
				Name: "test",
			},
			true,
			nil,
		},
		{
			expression.VarNameContext + `.name.endsWith('find') && ` + expression.VarNameContext + `.arguments.findIndex(a => a === "-exec") === -1`,
			command.CommandExecutionArguments{
				Name:      "find",
				Arguments: []string{"/"},
			},
			true,
			nil,
		},
		{
			expression.VarNameContext + `.name.endsWith('find') && ` + expression.VarNameContext + `.arguments.findIndex(a => a === "-exec") === -1`,
			command.CommandExecutionArguments{
				Name:      "find",
				Arguments: []string{"/", "-exec", "rm", "{}", ";"},
			},
			false,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			result, err := CalcApprovalExpr(tt.expression, tt.args)

			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
