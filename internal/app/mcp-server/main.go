package mcp_server

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rainu/ask-mai/internal/config"
	mcpServer "github.com/rainu/ask-mai/internal/mcp/server"
	"log/slog"
	"os"
)

type Args struct {
	VersionLine string
}

func Main(args Args) int {
	cfg := config.Parse(os.Args[1:], os.Environ())
	if cfg.Version {
		fmt.Fprintln(os.Stderr, args.VersionLine)
		return 0
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	slog.SetLogLoggerLevel(*cfg.DebugConfig.LogLevelParsed)

	ap := cfg.GetActiveProfile()

	ms := mcpServer.NewServer(args.VersionLine, ap.LLM.Tool.BuiltIns, ap.LLM.Tool.Custom)
	err := server.ServeStdio(ms)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 2
	}

	return 0
}
