package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileAppending struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileAppendingResult    `config:"-" yaml:"-"`
	Z file.FileAppendingArguments `config:"-" yaml:"-"`
}

func NewFileAppending() FileAppending {
	return FileAppending{
		Approval: ApprovalAlways,
	}
}

func (f FileAppending) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "appendFile",
		Approval:    f.Approval,
		Description: file.FileAppendingDefinition.Description,
		Parameters:  file.FileAppendingDefinition.Parameter,
		CommandFn:   file.FileAppendingDefinition.Function,
	}
}
