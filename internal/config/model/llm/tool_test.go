package llm

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	iMcp "github.com/rainu/ask-mai/internal/config/model/llm/mcp"
	"github.com/rainu/ask-mai/internal/expression"
	"github.com/rainu/ask-mai/internal/mcp/client"
	cmdchain "github.com/rainu/go-command-chain"
	"github.com/rainu/go-yacl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestTool_NeedApproval(t *testing.T) {
	tests := []struct {
		definition mcp.Tool
		approval   iMcp.Approval
		want       bool
	}{
		{
			approval: iMcp.Approval(""),
			want:     false,
		},
		{
			approval: iMcp.ApprovalNever,
			want:     false,
		},
		{
			approval: iMcp.ApprovalAlways,
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

func TestTool_GetTools(t *testing.T) {
	_, isCI := os.LookupEnv("CI")
	if isCI {
		t.Skip("Skipping test in CI environment")
		return
	}

	_, _, err := cmdchain.Builder().Join("docker", "-v").Finalize().RunAndGet()
	if err != nil {
		t.Skipf("Docker is not available: %v", err)
		return
	}

	defer client.Close()

	toTest := &LLMConfig{
		McpServer: map[string]iMcp.Server{
			"github": {
				Command: iMcp.Command{
					Name:      "docker",
					Arguments: []string{"run", "--rm", "-i", "-e", "GITHUB_PERSONAL_ACCESS_TOKEN=github_", "ghcr.io/github/github-mcp-server"},
				},
				Timeout: iMcp.Timeout{
					Init: yacl.P(1 * time.Second),
					List: yacl.P(1 * time.Second),
				},
			},
			"brave": {
				Command: iMcp.Command{
					Name:      "docker",
					Arguments: []string{"run", "--rm", "-i", "-e", "BRAVE_API_KEY=BS...", "mcp/brave-search"},
				},
				Timeout: iMcp.Timeout{
					Init: yacl.P(1 * time.Second),
					List: yacl.P(1 * time.Second),
				},
			},
		},
	}

	start := time.Now()
	result, err := toTest.GetTools(t.Context())
	duration := time.Since(start)

	require.NoError(t, err)
	assert.NotEmpty(t, result)

	var expectedMax time.Duration
	for _, server := range toTest.McpServer {
		expectedMax += *server.Timeout.Init
		expectedMax += *server.Timeout.List
	}
	expectedMax = expectedMax / time.Duration(len(toTest.McpServer))

	assert.True(t, duration < expectedMax, "took too long to merge tools - should be run in parallel")
}
