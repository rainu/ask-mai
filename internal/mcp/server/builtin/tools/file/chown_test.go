package file

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestTool_Chown(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeOwnerTool, ChangeOwnerToolHandler)
	})

	tf, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer func() {
		tf.Close()
		_ = os.Remove(tf.Name())
	}()

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeOwnerTool.Name
	req.Params.Arguments = map[string]any{
		"path":    tf.Name(),
		"user_id": 9,
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "operation not permitted")
}

func TestTool_Chown_Requirements(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeOwnerTool, ChangeOwnerToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeOwnerTool.Name
	req.Params.Arguments = map[string]any{
		"user_id": 9,
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
}

func TestTool_Chown_Requirements2(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(ChangeOwnerTool, ChangeOwnerToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = ChangeOwnerTool.Name
	req.Params.Arguments = map[string]any{
		"path": "/tmp",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'user_id' or 'group_id'")
}
