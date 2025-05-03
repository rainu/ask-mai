package mcp

import (
	"fmt"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/rainu/ask-mai/expression"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTool_NeedApproval(t *testing.T) {
	tests := []struct {
		definition mcp_golang.ToolRetType
		approval   Approval
		want       bool
	}{
		{
			approval: Approval(""),
			want:     false,
		},
		{
			approval: ApprovalNever,
			want:     false,
		},
		{
			approval: ApprovalAlways,
			want:     true,
		},
		{
			definition: mcp_golang.ToolRetType{Name: "test"},
			approval:   expression.VarNameContext + `.definition.name === "test"`,
			want:       true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestTool_NeedApproval_%d", i), func(t *testing.T) {
			tool := &Tool{
				ToolRetType: tt.definition,
				approval:    tt.approval,
			}

			assert.Equal(t, tt.want, tool.NeedApproval(t.Context(), "{}"))
		})
	}
}
