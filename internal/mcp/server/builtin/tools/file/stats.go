package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type StatsArguments struct {
	Path Path `json:"path"`
}

type StatsResult struct {
	Path        string    `json:"path"`
	IsDirectory bool      `json:"isDirectory"`
	IsRegular   bool      `json:"isRegular"`
	Permissions string    `json:"permissions"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"modTime"`
}

var StatsTool = mcp.NewTool("getStats",
	mcp.WithDescription("Get stats of a file or directory on the user's system."),
	mcp.WithString("path",
		mcp.Required(),
		mcp.Description("The path to the file or directory to get info for. Use '~' as placeholder for the user's home directory."),
	),
)

var StatsToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var pArgs StatsArguments

	r, w := io.Pipe()
	go func() {
		defer w.Close()

		json.NewEncoder(w).Encode(request.Params.Arguments)
	}()

	err := json.NewDecoder(r).Decode(&pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	if string(pArgs.Path) == "" {
		return nil, fmt.Errorf("missing parameter: 'path'")
	}
	path, err := pArgs.Path.Get()
	if err != nil {
		return nil, err
	}

	stats, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error getting stats: %w", err)
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		slog.Warn("Error getting absolute path!", "error", err)
		absolutePath = path
	}

	raw, err := json.Marshal(StatsResult{
		Path:        absolutePath,
		IsDirectory: stats.IsDir(),
		IsRegular:   stats.Mode().IsRegular(),
		Permissions: stats.Mode().Perm().String(),
		Size:        stats.Size(),
		ModTime:     stats.ModTime(),
	})
	return mcp.NewToolResultText(string(raw)), err
}
