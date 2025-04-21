package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileDeletion struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileDeletionResult    `config:"-" yaml:"-"`
	Z file.FileDeletionArguments `config:"-" yaml:"-"`
}

func NewFileDeletion() FileDeletion {
	return FileDeletion{
		Approval: ApprovalAlways,
	}
}

func (f FileDeletion) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "deleteFile",
		Approval:    f.Approval,
		Description: file.FileDeletionDefinition.Description,
		Parameters:  file.FileDeletionDefinition.Parameter,
		CommandFn:   file.FileDeletionDefinition.Function,
	}
}
