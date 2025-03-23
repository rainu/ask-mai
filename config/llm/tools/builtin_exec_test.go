package tools

import "testing"

func TestCalcApprovalExpr(t *testing.T) {
	tests := []struct {
		expression  string
		args        CommandExecutionArguments
		expected    bool
		expectedErr error
	}{
		{
			"v.Name == 'test'",
			CommandExecutionArguments{
				Name: "test",
			},
			true,
			nil,
		},
		{
			`v.Name.endsWith('find') && v.Arguments.findIndex(a => a === "-exec") === -1`,
			CommandExecutionArguments{
				Name:      "find",
				Arguments: []string{"/"},
			},
			true,
			nil,
		},
		{
			`v.Name.endsWith('find') && v.Arguments.findIndex(a => a === "-exec") === -1`,
			CommandExecutionArguments{
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
