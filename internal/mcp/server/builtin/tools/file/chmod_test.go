package file

import (
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestTool_Chmod(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeModeTool, ChangeModeToolHandler)
	})

	tf, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer func() {
		tf.Close()
		_ = os.Remove(tf.Name())
	}()

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeModeTool.Name
	req.Params.Arguments = map[string]any{
		"path":       tf.Name(),
		"permission": "0777",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	info, err := os.Stat(tf.Name())
	require.NoError(t, err)

	assert.Equal(t, os.FileMode(0777), info.Mode().Perm(), "File permissions should be set to 0777")
}

func TestTool_Chmod_Requirements(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeModeTool, ChangeModeToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeModeTool.Name
	req.Params.Arguments = map[string]any{
		"path": "/tmp",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing parameter: 'permission'")
	assert.Nil(t, res)
}

func TestTool_Chmod_Requirements2(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeModeTool, ChangeModeToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeModeTool.Name
	req.Params.Arguments = map[string]any{
		"permission": "0777",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
	assert.Nil(t, res)
}

func getTestClient(t *testing.T, serverConf func(s *server.MCPServer)) *client.Client {
	s := server.NewMCPServer(
		"ask-mai",
		"test-version",
		server.WithToolCapabilities(false),
	)
	serverConf(s)

	c := client.NewClient(transport.NewInProcessTransport(s))

	_, err := c.Initialize(t.Context(), mcp.InitializeRequest{})
	require.NoError(t, err)

	return c
}
