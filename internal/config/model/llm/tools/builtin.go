package tools

import (
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/rainu/ask-mai/internal/mcp/client"
	"github.com/rainu/ask-mai/internal/mcp/server"
	"github.com/rainu/ask-mai/internal/mcp/server/tools/command"
	"github.com/rainu/ask-mai/internal/mcp/server/tools/file"
	"github.com/rainu/ask-mai/internal/mcp/server/tools/http"
	"github.com/rainu/ask-mai/internal/mcp/server/tools/system"
	"reflect"
)

const BuiltInPrefix = "_"

type BuiltIns struct {
	SystemInfo  SystemInfo  `yaml:"system-info,omitempty" usage:"[System information] "`
	Environment Environment `yaml:"environment,omitempty" usage:"[Environment] "`
	SystemTime  SystemTime  `yaml:"system-time,omitempty" usage:"[System time] "`

	Stats Stats `yaml:"stats,omitempty" usage:"[Stats] "`

	ChangeMode  ChangeMode  `yaml:"change-mode,omitempty" usage:"[Change mode] "`
	ChangeOwner ChangeOwner `yaml:"change-owner,omitempty" usage:"[Change owner] "`
	ChangeTimes ChangeTimes `yaml:"change-times,omitempty" usage:"[Change times] "`

	FileCreation     FileCreation     `yaml:"file-creation,omitempty" usage:"[File creation] "`
	FileTempCreation FileTempCreation `yaml:"temp-file-creation,omitempty" usage:"[Temporary file creation] "`
	FileAppending    FileAppending    `yaml:"file-appending,omitempty" usage:"[File appending] "`
	FileReading      FileReading      `yaml:"file-reading,omitempty" usage:"[File reading] "`
	FileDeletion     FileDeletion     `yaml:"file-deletion,omitempty,omitempty" usage:"[File deletion] "`

	DirectoryCreation     DirectoryCreation     `yaml:"dir-creation,omitempty" usage:"[Directory creation] "`
	DirectoryTempCreation DirectoryTempCreation `yaml:"temp-dir-creation,omitempty" usage:"[Temporary directory creation] "`
	DirectoryDeletion     DirectoryDeletion     `yaml:"dir-deletion,omitempty" usage:"[Directory deletion] "`

	CommandExec CommandExecution `yaml:"command-execution,omitempty" usage:"[Command execution] "`

	Http Http `yaml:"http,omitempty" usage:"[HTTP] "`

	Disable bool `yaml:"disable,omitempty" usage:"Disable all builtin tools."`
}

func (b *BuiltIns) GetTransport() (transport.Interface, error) {
	return transport.NewInProcessTransport(server.NewServer("", server.Options{
		DisableSystemInfo:            b.SystemInfo.Disable,
		DisableEnvironment:           b.Environment.Disable,
		DisableSystemTime:            b.SystemTime.Disable,
		DisableStats:                 b.Stats.Disable,
		DisableChangeMode:            b.ChangeMode.Disable,
		DisableChangeOwner:           b.ChangeOwner.Disable,
		DisableChangeTimes:           b.ChangeTimes.Disable,
		DisableFileCreation:          b.FileCreation.Disable,
		DisableFileTempCreation:      b.FileTempCreation.Disable,
		DisableFileAppending:         b.FileAppending.Disable,
		DisableFileReading:           b.FileReading.Disable,
		DisableFileDeletion:          b.FileDeletion.Disable,
		DisableDirectoryCreation:     b.DirectoryCreation.Disable,
		DisableDirectoryTempCreation: b.DirectoryTempCreation.Disable,
		DisableDirectoryDeletion:     b.DirectoryDeletion.Disable,
		DisableCommandExec:           b.CommandExec.Disable,
		DisableHttp:                  b.Http.Disable,
	})), nil
}

func (b *BuiltIns) GetTimeouts() client.Timeouts {
	return client.Timeouts{} // No timeouts for built-in tools
}

func (b *BuiltIns) GetApprovalFor(toolName string) string {
	if system.SystemInfoTool.Name == toolName {
		return b.SystemInfo.Approval
	} else if system.EnvironmentTool.Name == toolName {
		return b.Environment.Approval
	} else if system.SystemTimeTool.Name == toolName {
		return b.SystemTime.Approval
	} else if file.StatsTool.Name == toolName {
		return b.Stats.Approval
	} else if file.ChangeModeTool.Name == toolName {
		return b.ChangeMode.Approval
	} else if file.ChangeOwnerTool.Name == toolName {
		return b.ChangeOwner.Approval
	} else if file.ChangeTimesTool.Name == toolName {
		return b.ChangeTimes.Approval
	} else if file.FileCreationTool.Name == toolName {
		return b.FileCreation.Approval
	} else if file.FileTempCreationTool.Name == toolName {
		return b.FileTempCreation.Approval
	} else if file.FileAppendingTool.Name == toolName {
		return b.FileAppending.Approval
	} else if file.FileReadingTool.Name == toolName {
		return b.FileReading.Approval
	} else if file.FileDeletionTool.Name == toolName {
		return b.FileDeletion.Approval
	} else if file.DirectoryCreationTool.Name == toolName {
		return b.DirectoryCreation.Approval
	} else if file.DirectoryTempCreationTool.Name == toolName {
		return b.DirectoryTempCreation.Approval
	} else if file.DirectoryDeletionTool.Name == toolName {
		return b.DirectoryDeletion.Approval
	} else if command.CommandExecutionTool.Name == toolName {
		return b.CommandExec.Approval
	} else if http.CallTool.Name == toolName {
		return b.Http.Approval
	}

	return ApprovalAlways
}

func NewBuiltIns() *BuiltIns {
	return &BuiltIns{
		SystemInfo:            NewSystemInfo(),
		Environment:           NewEnvironment(),
		SystemTime:            NewSystemTime(),
		Stats:                 NewStats(),
		ChangeMode:            NewChangeMode(),
		ChangeOwner:           NewChangeOwner(),
		ChangeTimes:           NewChangeTimes(),
		FileCreation:          NewFileCreation(),
		FileTempCreation:      NewFileTempCreation(),
		FileAppending:         NewFileAppending(),
		FileReading:           NewFileReading(),
		FileDeletion:          NewFileDeletion(),
		DirectoryCreation:     NewDirectoryCreation(),
		DirectoryTempCreation: NewDirectoryTempCreation(),
		DirectoryDeletion:     NewDirectoryDeletion(),
		CommandExec:           NewCommandExecution(),
		Http:                  NewHttp(),
		Disable:               false,
	}
}

func (b BuiltIns) AsFunctionDefinitions() []FunctionDefinition {
	if b.Disable {
		return []FunctionDefinition{}
	}

	var functions []FunctionDefinition

	v := reflect.ValueOf(b)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Check if the field has a method called AsFunctionDefinition
		method := field.MethodByName("AsFunctionDefinition")
		if !method.IsValid() {
			continue
		}

		// Call the method and get the result
		result := method.Call(nil)
		if len(result) == 0 || result[0].IsNil() {
			continue
		}

		// Convert the result to FunctionDefinition and add it to the list
		if fn, ok := result[0].Interface().(*FunctionDefinition); ok && fn != nil {
			fn.Name = BuiltInPrefix + fn.Name
			functions = append(functions, *fn)
		}
	}

	return functions
}
