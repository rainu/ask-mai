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

func TestTool_FileAppending(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileAppendingTool, FileAppendingToolHandler)
	})

	testFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	_, err = testFile.WriteString("First line.\n")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileAppendingTool.Name
	req.Params.Arguments = map[string]any{
		"path":    testFile.Name(),
		"content": "Second line.\n",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testFile.Name())

	// re-open the file to verify the content
	testFile, err = os.Open(testFile.Name())
	require.NoError(t, err)

	rawContent, err := io.ReadAll(testFile)
	assert.NoError(t, err)
	assert.Equal(t, "First line.\nSecond line.\n", string(rawContent))
}

func TestTool_FileAppending_HomeResolving(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileAppendingTool, FileAppendingToolHandler)
	})

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	testFile, err := os.Create(path.Join(home, fmt.Sprintf(".ask-mai-%d", time.Now().Unix())))
	require.NoError(t, err)
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	_, err = testFile.WriteString("First line.\n")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileAppendingTool.Name
	req.Params.Arguments = map[string]any{
		"path":    path.Join("~", path.Base(testFile.Name())),
		"content": "Second line.\n",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, testFile.Name())

	// re-open the file to verify the content
	testFile, err = os.Open(testFile.Name())
	require.NoError(t, err)

	rawContent, err := io.ReadAll(testFile)
	assert.NoError(t, err)
	assert.Equal(t, "First line.\nSecond line.\n", string(rawContent))
}

func TestTool_FileAppending_Requirements_Path(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileAppendingTool, FileAppendingToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileAppendingTool.Name
	req.Params.Arguments = map[string]any{}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
}

func TestTool_FileAppending_Requirements_Content(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileAppendingTool, FileAppendingToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileAppendingTool.Name
	req.Params.Arguments = map[string]any{
		"path": "some/path/to/file.txt",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'content'")
}

func TestTool_FileAppending_Requirements_FileDoesNotExists(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileAppendingTool, FileAppendingToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileAppendingTool.Name
	req.Params.Arguments = map[string]any{
		"path":    "some/path/to/file.txt",
		"content": "This is a test content.",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "file does not exists")
}

func TestTool_FileAppending_Requirements_PathIsDirectory(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileAppendingTool, FileAppendingToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileAppendingTool.Name
	req.Params.Arguments = map[string]any{
		"path":    t.TempDir(),
		"content": "This is a test content.",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "path is a directory")
}
