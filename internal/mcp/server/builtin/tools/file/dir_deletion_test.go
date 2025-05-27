package file

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestTool_DirDeletion(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(DirectoryDeletionTool, DirectoryDeletionToolHandler)
	})

	testDir := path.Join(t.TempDir(), t.Name(), "testdir")
	require.NoError(t, os.MkdirAll(testDir, os.ModePerm))

	req := mcp.CallToolRequest{}
	req.Params.Name = DirectoryDeletionTool.Name
	req.Params.Arguments = map[string]any{
		"path": testDir,
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testDir)

	_, err = os.Stat(testDir)
	assert.Error(t, err)
}

func TestTool_DirDeletion_Requirements(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(DirectoryDeletionTool, DirectoryDeletionToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = DirectoryDeletionTool.Name
	req.Params.Arguments = map[string]any{}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
}
