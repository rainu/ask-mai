package file

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path"
	"testing"
	"time"
)

func TestTool_FileCreation(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileCreationTool, FileCreationToolHandler)
	})

	testFilePath := path.Join(t.TempDir(), "test.txt")

	req := mcp.CallToolRequest{}
	req.Params.Name = FileCreationTool.Name
	req.Params.Arguments = map[string]any{
		"path":    testFilePath,
		"content": "Hello World!\n",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testFilePath)

	// open the file to verify the content
	testFile, err := os.Open(testFilePath)
	require.NoError(t, err)

	rawContent, err := io.ReadAll(testFile)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!\n", string(rawContent))

	info, err := os.Stat(testFile.Name())
	require.NoError(t, err)

	assert.Equal(t, os.FileMode(0644), info.Mode().Perm(), "File permissions should be set to 0644")
}

func TestTool_FileCreation_CustomPermissions(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileCreationTool, FileCreationToolHandler)
	})

	testFilePath := path.Join(t.TempDir(), "test.txt")

	req := mcp.CallToolRequest{}
	req.Params.Name = FileCreationTool.Name
	req.Params.Arguments = map[string]any{
		"path":       testFilePath,
		"permission": "0600",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testFilePath)

	info, err := os.Stat(testFilePath)
	require.NoError(t, err)

	assert.Equal(t, os.FileMode(0600), info.Mode().Perm(), "File permissions should be set to 0600")
}

func TestTool_FileCreation_HomeResolving(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileCreationTool, FileCreationToolHandler)
	})

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	testFilePath := path.Join(home, fmt.Sprintf(".ask-mai-%d", time.Now().Unix()))
	require.NoError(t, err)
	defer func() {
		os.Remove(testFilePath)
	}()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileCreationTool.Name
	req.Params.Arguments = map[string]any{
		"path":    path.Join("~", path.Base(testFilePath)),
		"content": "Hello World!\n",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testFilePath)

	// open the file to verify the content
	testFile, err := os.Open(testFilePath)
	require.NoError(t, err)

	rawContent, err := io.ReadAll(testFile)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!\n", string(rawContent))

	info, err := os.Stat(testFile.Name())
	require.NoError(t, err)

	assert.Equal(t, os.FileMode(0644), info.Mode().Perm(), "File permissions should be set to 0644")
}

func TestTool_FileCreation_Requirements_Path(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileCreationTool, FileCreationToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileCreationTool.Name
	req.Params.Arguments = map[string]any{}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
}

func TestTool_FileCreation_Requirements_FileDoesExists(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileCreationTool, FileCreationToolHandler)
	})
	temp, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	req := mcp.CallToolRequest{}
	req.Params.Name = FileCreationTool.Name
	req.Params.Arguments = map[string]any{
		"path":    temp.Name(),
		"content": "This is a test content.",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "file already exists")
}

func TestTool_FileCreation_Requirements_PathIsDirectory(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileCreationTool, FileCreationToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileCreationTool.Name
	req.Params.Arguments = map[string]any{
		"path":    t.TempDir(),
		"content": "This is a test content.",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "path exists but is a directory")
}
