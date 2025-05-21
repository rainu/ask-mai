package server

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rainu/ask-mai/internal/mcp/server/tools/command"
	"github.com/rainu/ask-mai/internal/mcp/server/tools/file"
	"github.com/rainu/ask-mai/internal/mcp/server/tools/http"
	"github.com/rainu/ask-mai/internal/mcp/server/tools/system"
)

type Options struct {
	DisableSystemInfo            bool
	DisableEnvironment           bool
	DisableSystemTime            bool
	DisableStats                 bool
	DisableChangeMode            bool
	DisableChangeOwner           bool
	DisableChangeTimes           bool
	DisableFileCreation          bool
	DisableFileTempCreation      bool
	DisableFileAppending         bool
	DisableFileReading           bool
	DisableFileDeletion          bool
	DisableDirectoryCreation     bool
	DisableDirectoryTempCreation bool
	DisableDirectoryDeletion     bool
	DisableCommandExec           bool
	DisableHttp                  bool
}

func NewServer(version string, opts Options) *server.MCPServer {
	s := server.NewMCPServer(
		"ask-mai",
		version,
		server.WithToolCapabilities(false),
	)

	if !opts.DisableSystemTime {
		s.AddTool(system.SystemTimeTool, system.SystemTimeToolHandler)
	}
	if !opts.DisableSystemInfo {
		s.AddTool(system.SystemInfoTool, system.SystemInfoToolHandler)
	}
	if !opts.DisableEnvironment {
		s.AddTool(system.EnvironmentTool, system.EnvironmentToolHandler)
	}

	if !opts.DisableChangeMode {
		s.AddTool(file.ChangeModeTool, file.ChangeModeToolHandler)
	}
	if !opts.DisableChangeOwner {
		s.AddTool(file.ChangeOwnerTool, file.ChangeOwnerToolHandler)
	}
	if !opts.DisableChangeTimes {
		s.AddTool(file.ChangeTimesTool, file.ChangeTimesToolHandler)
	}

	if !opts.DisableDirectoryCreation {
		s.AddTool(file.DirectoryCreationTool, file.DirectoryCreationToolHandler)
	}
	if !opts.DisableDirectoryDeletion {
		s.AddTool(file.DirectoryDeletionTool, file.DirectoryDeletionToolHandler)
	}
	if !opts.DisableDirectoryTempCreation {
		s.AddTool(file.DirectoryTempCreationTool, file.DirectoryTempCreationToolHandler)
	}

	if !opts.DisableFileAppending {
		s.AddTool(file.FileAppendingTool, file.FileAppendingToolHandler)
	}
	if !opts.DisableFileCreation {
		s.AddTool(file.FileCreationTool, file.FileCreationToolHandler)
	}
	if !opts.DisableFileDeletion {
		s.AddTool(file.FileDeletionTool, file.FileDeletionToolHandler)
	}
	if !opts.DisableFileReading {
		s.AddTool(file.FileReadingTool, file.FileReadingToolHandler)
	}
	if !opts.DisableFileTempCreation {
		s.AddTool(file.FileTempCreationTool, file.FileTempCreationToolHandler)
	}
	if !opts.DisableStats {
		s.AddTool(file.StatsTool, file.StatsToolHandler)
	}

	if !opts.DisableCommandExec {
		s.AddTool(command.CommandExecutionTool, command.CommandExecutionToolHandler)
	}

	if !opts.DisableHttp {
		s.AddTool(http.CallTool, http.CallToolHandler)
	}
	s.DeleteTools()

	return s
}
