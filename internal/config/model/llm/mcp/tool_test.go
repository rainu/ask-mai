package mcp

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/expression"
	internalMcp "github.com/rainu/ask-mai/internal/llms/tools/mcp"
	cmdchain "github.com/rainu/go-command-chain"
	"github.com/rainu/go-yacl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTool_NeedApproval(t *testing.T) {
	tests := []struct {
		definition mcp.Tool
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
			definition: mcp.Tool{Name: "test"},
			approval:   expression.VarNameContext + `.definition.name === "test"`,
			want:       true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestTool_NeedApproval_%d", i), func(t *testing.T) {
			tool := &Tool{
				Tool:     tt.definition,
				approval: tt.approval,
			}

			assert.Equal(t, tt.want, tool.NeedApproval(t.Context(), "{}"))
		})
	}
}

func TestMergeTools(t *testing.T) {
	_, _, err := cmdchain.Builder().Join("docker", "-v").Finalize().RunAndGet()
	if err != nil {
		t.Skipf("Docker is not available: %v", err)
		return
	}

	defer internalMcp.Close()

	toTest := map[string]Server{
		"github": {
			Command: Command{
				Name:      "docker",
				Arguments: []string{"run", "--rm", "-i", "-e", "GITHUB_PERSONAL_ACCESS_TOKEN=github_", "ghcr.io/github/github-mcp-server"},
			},
			Timeout: Timeout{
				Init: yacl.P(1 * time.Second),
				List: yacl.P(1 * time.Second),
			},
		},
		"brave": {
			Command: Command{
				Name:      "docker",
				Arguments: []string{"run", "--rm", "-i", "-e", "BRAVE_API_KEY=BS...", "mcp/brave-search"},
			},
			Timeout: Timeout{
				Init: yacl.P(1 * time.Second),
				List: yacl.P(1 * time.Second),
			},
		},
	}

	start := time.Now()
	result, err := MergeTools(t.Context(), toTest)
	duration := time.Since(start)

	require.NoError(t, err)
	assert.NotEmpty(t, result)

	var expectedMax time.Duration
	for _, server := range toTest {
		expectedMax += *server.Timeout.Init
		expectedMax += *server.Timeout.List
	}
	expectedMax = expectedMax / time.Duration(len(toTest))

	assert.True(t, duration < expectedMax, "took too long to merge tools - should be run in parallel")
}
