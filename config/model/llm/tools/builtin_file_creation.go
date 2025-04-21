package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileCreation struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileCreationResult    `config:"-" yaml:"-"`
	Z file.FileCreationArguments `config:"-" yaml:"-"`
}

func NewFileCreation() FileCreation {
	return FileCreation{
		Approval: ApprovalNever,
	}
}

func (f FileCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "createFile",
		Approval:    f.Approval,
		Description: file.FileCreationDefinition.Description,
		Parameters:  file.FileCreationDefinition.Parameter,
		CommandFn:   file.FileCreationDefinition.Function,
	}
}
