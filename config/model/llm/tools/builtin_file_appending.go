package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileAppending struct {
	Disable    bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval bool `config:"no-approval" yaml:"no-approval" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.FileAppendingResult    `config:"-" yaml:"-"`
	Z file.FileAppendingArguments `config:"-" yaml:"-"`
}

func (f FileAppending) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "appendFile",
		NeedsApproval: !f.NoApproval,
		Description:   file.FileAppendingDefinition.Description,
		Parameters:    file.FileAppendingDefinition.Parameter,
		CommandFn:     file.FileAppendingDefinition.Function,
	}
}
