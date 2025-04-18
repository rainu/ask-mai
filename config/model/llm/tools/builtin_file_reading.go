package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileReading struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.FileReadingResult    `config:"-" yaml:"-"`
	Z file.FileReadingArguments `config:"-" yaml:"-"`
}

func (f FileReading) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "readTextFile",
		NeedsApproval: f.NeedsApproval,
		Description:   file.FileReadingDefinition.Description,
		Parameters:    file.FileReadingDefinition.Parameter,
		CommandFn:     file.FileReadingDefinition.Function,
	}
}
