package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileTempCreation struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.FileTempCreationResult    `config:"-" yaml:"-"`
	Z file.FileTempCreationArguments `config:"-" yaml:"-"`
}

func (f FileTempCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "createTempFile",
		NeedsApproval: f.NeedsApproval,
		Description:   file.FileTempCreationDefinition.Description,
		Parameters:    file.FileTempCreationDefinition.Parameter,
		CommandFn:     file.FileTempCreationDefinition.Function,
	}
}
