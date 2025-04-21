package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type DirectoryTempCreation struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.DirectoryTempCreationResult    `config:"-" yaml:"-"`
	Z file.DirectoryTempCreationArguments `config:"-" yaml:"-"`
}

func NewDirectoryTempCreation() DirectoryTempCreation {
	return DirectoryTempCreation{
		Approval: ApprovalNever,
	}
}

func (f DirectoryTempCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "createTempDirectory",
		Approval:    f.Approval,
		Description: file.DirectoryTempCreationDefinition.Description,
		Parameters:  file.DirectoryTempCreationDefinition.Parameter,
		CommandFn:   file.DirectoryTempCreationDefinition.Function,
	}
}
