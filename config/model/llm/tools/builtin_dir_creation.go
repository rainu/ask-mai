package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type DirectoryCreation struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.DirectoryCreationResult    `config:"-" yaml:"-"`
	Z file.DirectoryCreationArguments `config:"-" yaml:"-"`
}

func NewDirectoryCreation() DirectoryCreation {
	return DirectoryCreation{
		Approval: ApprovalNever,
	}
}

func (f DirectoryCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "createDirectory",
		Approval:    f.Approval,
		Description: file.DirectoryCreationDefinition.Description,
		Parameters:  file.DirectoryCreationDefinition.Parameter,
		CommandFn:   file.DirectoryCreationDefinition.Function,
	}
}
