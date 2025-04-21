package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileReading struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileReadingResult    `config:"-" yaml:"-"`
	Z file.FileReadingArguments `config:"-" yaml:"-"`
}

func NewFileReading() FileReading {
	return FileReading{
		Approval: ApprovalNever,
	}
}

func (f FileReading) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "readTextFile",
		Approval:    f.Approval,
		Description: file.FileReadingDefinition.Description,
		Parameters:  file.FileReadingDefinition.Parameter,
		CommandFn:   file.FileReadingDefinition.Function,
	}
}
