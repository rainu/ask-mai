package system

import (
	"context"
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"runtime"
)

type SystemInfoArguments struct {
}

type SystemInfoResult struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	CPU      int    `json:"cpus"`
	Hostname string `json:"hostname"`
	UserDir  string `json:"user_dir"`
	UserId   int    `json:"user_id"`
	GroupId  int    `json:"group_id"`
	WorkDir  string `json:"working_directory"`
	PID      int    `json:"process_id"`
}

var SystemInfoTool = mcp.NewTool("getSystemInformation",
	mcp.WithDescription("Get the following information about the user's system: OS, architecture, number of CPUs, hostname, user directory, user ID, group ID, working directory, process ID."),
)

var SystemInfoToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	raw, err := json.Marshal(SystemInfoResult{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		CPU:  runtime.NumCPU(),
		Hostname: func() string {
			h, err := os.Hostname()
			if err != nil {
				return "unknown"
			}
			return h
		}(),
		UserDir: func() string {
			home, err := os.UserHomeDir()
			if err != nil {
				return "unknown"
			}
			return home
		}(),
		UserId:  os.Getuid(),
		GroupId: os.Getgid(),
		WorkDir: func() string {
			dir, err := os.Getwd()
			if err != nil {
				return "unknown"
			}
			return dir
		}(),
		PID: os.Getpid(),
	})
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(raw)), nil
}
