package file

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"time"
)

func TestTool_DirCreation(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(DirectoryCreationTool, DirectoryCreationToolHandler)
	})

	testDir := path.Join(t.TempDir(), t.Name(), "testdir")

	req := mcp.CallToolRequest{}
	req.Params.Name = DirectoryCreationTool.Name
	req.Params.Arguments = map[string]any{
		"path": testDir,
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testDir)

	info, err := os.Stat(testDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestTool_DirCreation_WithHome(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(DirectoryCreationTool, DirectoryCreationToolHandler)
	})
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	testDirName := fmt.Sprintf(".ask-mai-%d", time.Now().Unix())
	testDirPath := path.Join(home, testDirName)
	defer func() {
		os.Remove(testDirPath)
	}()

	req := mcp.CallToolRequest{}
	req.Params.Name = DirectoryCreationTool.Name
	req.Params.Arguments = map[string]any{
		"path": path.Join("~", testDirName),
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testDirPath)

	info, err := os.Stat(testDirPath)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestTool_DirCreation_WithPermissions(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(DirectoryCreationTool, DirectoryCreationToolHandler)
	})

	testDir := path.Join(t.TempDir(), t.Name(), "testdir")

	req := mcp.CallToolRequest{}
	req.Params.Name = DirectoryCreationTool.Name
	req.Params.Arguments = map[string]any{
		"path":       testDir,
		"permission": "0700",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testDir)

	info, err := os.Stat(testDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
	assert.Equal(t, os.FileMode(0700), info.Mode().Perm(), "Directory permissions should be set to 0700")
}

func TestTool_DirCreation_Requirements(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(DirectoryCreationTool, DirectoryCreationToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = DirectoryCreationTool.Name
	req.Params.Arguments = map[string]any{}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
}
