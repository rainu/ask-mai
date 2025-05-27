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

func TestTool_FileDeletion(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileDeletionTool, FileDeletionToolHandler)
	})
	testFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileDeletionTool.Name
	req.Params.Arguments = map[string]any{
		"path": testFile.Name(),
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testFile.Name())

	_, err = os.Stat(testFile.Name())
	assert.True(t, os.IsNotExist(err))
}

func TestTool_FileDeletion_HomeResolving(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileDeletionTool, FileDeletionToolHandler)
	})

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	testFilePath := path.Join(home, fmt.Sprintf(".ask-mai-%d", time.Now().Unix()))
	require.NoError(t, err)
	defer func() {
		os.Remove(testFilePath)
	}()

	testFile, err := os.Create(testFilePath)
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileDeletionTool.Name
	req.Params.Arguments = map[string]any{
		"path": path.Join("~", path.Base(testFilePath)),
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testFilePath)

	_, err = os.Stat(testFile.Name())
	assert.True(t, os.IsNotExist(err))
}

func TestTool_FileDeletion_Requirements_Path(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileDeletionTool, FileDeletionToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileDeletionTool.Name
	req.Params.Arguments = map[string]any{}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
}

func TestTool_FileDeletion_Requirements_FileDoesNotExists(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileDeletionTool, FileDeletionToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileDeletionTool.Name
	req.Params.Arguments = map[string]any{
		"path": "/file/does/not/exist.txt",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "file does not exist")
}

func TestTool_FileDeletion_Requirements_PathIsDirectory(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileDeletionTool, FileDeletionToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileDeletionTool.Name
	req.Params.Arguments = map[string]any{
		"path": t.TempDir(),
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "path is a directory")
}
