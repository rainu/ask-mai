package file

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestTool_Chtime(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeTimesTool, ChangeTimesToolHandler)
	})

	tf, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer func() {
		tf.Close()
		_ = os.Remove(tf.Name())
	}()

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeTimesTool.Name
	req.Params.Arguments = map[string]any{
		"path":              tf.Name(),
		"modification_time": "2023-01-01T00:00:00Z",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	info, err := os.Stat(tf.Name())
	require.NoError(t, err)
	assert.Contains(t, info.ModTime().String(), "2023-01-01")
}

func TestTool_Chtime_InvalidTime_access(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeTimesTool, ChangeTimesToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeTimesTool.Name
	req.Params.Arguments = map[string]any{
		"path":        "/tmp",
		"access_time": "invalid-time",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "cannot parse")
}

func TestTool_Chtime_InvalidTime_modification(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeTimesTool, ChangeTimesToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeTimesTool.Name
	req.Params.Arguments = map[string]any{
		"path":              "/tmp",
		"modification_time": "invalid-time",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "cannot parse")
}

func TestTool_Chtime_Requirements(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeTimesTool, ChangeTimesToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeTimesTool.Name
	req.Params.Arguments = map[string]any{}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
}

func TestTool_Chtime_Requirements2(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeTimesTool, ChangeTimesToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeTimesTool.Name
	req.Params.Arguments = map[string]any{
		"path": "/tmp",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: at least one of 'access_time' or 'modification_time' must be set")
}
