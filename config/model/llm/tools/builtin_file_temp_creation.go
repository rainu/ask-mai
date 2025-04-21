package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileTempCreation struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileTempCreationResult    `config:"-" yaml:"-"`
	Z file.FileTempCreationArguments `config:"-" yaml:"-"`
}

func NewFileTempCreation() FileTempCreation {
	return FileTempCreation{
		Approval: ApprovalNever,
	}
}

func (f FileTempCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "createTempFile",
		Approval:    f.Approval,
		Description: file.FileTempCreationDefinition.Description,
		Parameters:  file.FileTempCreationDefinition.Parameter,
		CommandFn:   file.FileTempCreationDefinition.Function,
	}
}
