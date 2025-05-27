package command

import (
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTool_Command_Exec_Echo(t *testing.T) {
	c := getTestClient(t)

	req := mcp.CallToolRequest{}
	req.Params.Name = CommandExecutionTool.Name
	req.Params.Arguments = map[string]any{
		"name":      "echo",
		"arguments": []string{"hello", "world"},
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	text := res.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "hello world")
}

func TestTool_Command_Exec_Env(t *testing.T) {
	c := getTestClient(t)

	req := mcp.CallToolRequest{}
	req.Params.Name = CommandExecutionTool.Name
	req.Params.Arguments = map[string]any{
		"name": "env",
		"environment": map[string]string{
			"FOO": "bar",
		},
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	text := res.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "FOO=bar")
}

func TestTool_Command_Exec_Unknown(t *testing.T) {
	c := getTestClient(t)

	req := mcp.CallToolRequest{}
	req.Params.Name = CommandExecutionTool.Name
	req.Params.Arguments = map[string]any{
		"name": "CommandShouldNotExists",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `failed to start command: exec: "CommandShouldNotExists": executable file not found`)
	assert.Nil(t, res)
}

func getTestClient(t *testing.T) *client.Client {
	s := server.NewMCPServer(
		"ask-mai",
		"test-version",
		server.WithToolCapabilities(false),
	)
	s.AddTool(CommandExecutionTool, CommandExecutionToolHandler)

	c := client.NewClient(transport.NewInProcessTransport(s))

	_, err := c.Initialize(t.Context(), mcp.InitializeRequest{})
	require.NoError(t, err)

	return c
}
