package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileCreation struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"Disable tool"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileCreationResult    `yaml:"-"`
	Z file.FileCreationArguments `yaml:"-"`
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
