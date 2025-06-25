package command

import (
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
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

func TestTool_Command_Exec_Sudo(t *testing.T) {
	_, isCI := os.LookupEnv("CI")
	if isCI {
		t.Skip("Skipping test in CI environment")
		return
	}

	c := getTestClient(t)

	apScript, err := os.Create(path.Join(t.TempDir(), "ask_pass.sh"))
	require.NoError(t, err)
	require.NoError(t, os.Chmod(apScript.Name(), 0700)) // Make it executable

	_, err = apScript.WriteString(`#!/bin/sh
touch ` + path.Dir(apScript.Name()) + `/called
echo 'invalidPassword'
`)
	require.NoError(t, err)
	require.NoError(t, apScript.Close())

	os.Setenv("SUDO_ASKPASS", apScript.Name())

	req := mcp.CallToolRequest{}
	req.Params.Name = CommandExecutionTool.Name
	req.Params.Arguments = map[string]any{
		"name":      "sudo",
		"arguments": []string{"echo", "hello", "world"},
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// Check if the script was called
	calledFile := path.Join(path.Dir(apScript.Name()), "called")
	_, err = os.Stat(calledFile)
	assert.NoError(t, err, "The sudo askpass script should have been called")
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
