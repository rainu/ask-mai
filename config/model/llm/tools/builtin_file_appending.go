package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileAppending struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"Disable tool"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileAppendingResult    `yaml:"-"`
	Z file.FileAppendingArguments `yaml:"-"`
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
