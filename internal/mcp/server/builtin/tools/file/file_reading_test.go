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

func TestTool_FileReading(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
	})

	testFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	_, err = testFile.WriteString("First line.\nSecond line.\n")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{
		"path": testFile.Name(),
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "First line.\\nSecond line.\\n")
}

func TestTool_FileReading_LineLimit(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
	})

	testFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	_, err = testFile.WriteString("First line.\nSecond line.\n")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{
		"path": testFile.Name(),
		"lm":   "line",
		"ll":   1,
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "First line.")
	assert.NotContains(t, res.Content[0].(mcp.TextContent).Text, "Second line.")
}

func TestTool_FileReading_LineLimitWithOffset(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
	})

	testFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	_, err = testFile.WriteString("First line.\nSecond line.\n")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{
		"path": testFile.Name(),
		"lm":   "line",
		"lo":   1,
		"ll":   1,
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Second line.")
	assert.NotContains(t, res.Content[0].(mcp.TextContent).Text, "First line.")
}

func TestTool_FileReading_CharLimit(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
	})

	testFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	_, err = testFile.WriteString("First line.\nSecond line.\n")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{
		"path": testFile.Name(),
		"lm":   "char",
		"ll":   5,
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "First")
	assert.NotContains(t, res.Content[0].(mcp.TextContent).Text, "Seco")
}

func TestTool_FileReading_CharLimitWithOffset(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
	})

	testFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	_, err = testFile.WriteString("First line.\nSecond line.\n")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{
		"path": testFile.Name(),
		"lm":   "char",
		"lo":   12,
		"ll":   6,
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "Second")
	assert.NotContains(t, res.Content[0].(mcp.TextContent).Text, "First")
}

func TestTool_FileReading_HomeResolving(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
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
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{
		"path": path.Join("~", path.Base(testFile.Name())),
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "First line.\\n")
}

func TestTool_FileReading_Requirements_Path(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
}

func TestTool_FileReading_Requirements_InvalidLimitMode(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{
		"path": "some/path/to/file.txt",
		"lm":   "unknown",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "invalid limit mode")
}

func TestTool_FileReading_Requirements_FileDoesNotExists(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{
		"path": "some/path/to/file.txt",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "no such file")
}

func TestTool_FileReading_Requirements_PathIsDirectory(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileReadingTool, FileReadingToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileReadingTool.Name
	req.Params.Arguments = map[string]any{
		"path":    t.TempDir(),
		"content": "This is a test content.",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "is a directory")
}
