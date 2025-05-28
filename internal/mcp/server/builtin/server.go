package builtin

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/builtin"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/command"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/file"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/http"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/system"
)

func NewServer(version string, cfg builtin.BuiltIns) *server.MCPServer {
	s := server.NewMCPServer(
		"ask-mai",
		version,
		server.WithToolCapabilities(false),
	)
	AddTools(s, cfg)

	return s
}

func AddTools(s *server.MCPServer, cfg builtin.BuiltIns) {
	if !cfg.SystemTime.Disable {
		s.AddTool(system.SystemTimeTool, system.SystemTimeToolHandler)
	}
	if !cfg.SystemInfo.Disable {
		s.AddTool(system.SystemInfoTool, system.SystemInfoToolHandler)
	}
	if !cfg.Environment.Disable {
		s.AddTool(system.EnvironmentTool, system.EnvironmentToolHandler)
	}

	if !cfg.ChangeMode.Disable {
		s.AddTool(file.ChangeModeTool, file.ChangeModeToolHandler)
	}
	if !cfg.ChangeOwner.Disable {
		s.AddTool(file.ChangeOwnerTool, file.ChangeOwnerToolHandler)
	}
	if !cfg.ChangeTimes.Disable {
		s.AddTool(file.ChangeTimesTool, file.ChangeTimesToolHandler)
	}

	if !cfg.DirectoryCreation.Disable {
		s.AddTool(file.DirectoryCreationTool, file.DirectoryCreationToolHandler)
	}
	if !cfg.DirectoryDeletion.Disable {
		s.AddTool(file.DirectoryDeletionTool, file.DirectoryDeletionToolHandler)
	}
	if !cfg.DirectoryTempCreation.Disable {
		s.AddTool(file.DirectoryTempCreationTool, file.DirectoryTempCreationToolHandler)
	}

	if !cfg.FileAppending.Disable {
		s.AddTool(file.FileAppendingTool, file.FileAppendingToolHandler)
	}
	if !cfg.FileCreation.Disable {
		s.AddTool(file.FileCreationTool, file.FileCreationToolHandler)
	}
	if !cfg.FileDeletion.Disable {
		s.AddTool(file.FileDeletionTool, file.FileDeletionToolHandler)
	}
	if !cfg.FileReading.Disable {
		s.AddTool(file.FileReadingTool, file.FileReadingToolHandler)
	}
	if !cfg.FileTempCreation.Disable {
		s.AddTool(file.FileTempCreationTool, file.FileTempCreationToolHandler)
	}
	if !cfg.Stats.Disable {
		s.AddTool(file.StatsTool, file.StatsToolHandler)
	}

	if !cfg.CommandExec.Disable {
		s.AddTool(command.CommandExecutionTool, command.CommandExecutionToolHandler)
	}

	if !cfg.Http.Disable {
		s.AddTool(http.CallTool, http.CallToolHandler)
	}
}
