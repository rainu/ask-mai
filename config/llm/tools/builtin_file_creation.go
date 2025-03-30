package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileCreation struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.FileCreationResult    `config:"-" yaml:"-"`
	Z file.FileCreationArguments `config:"-" yaml:"-"`
}

func (f FileCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "createFile",
		NeedsApproval: f.NeedsApproval,
		Description:   file.FileCreationDefinition.Description,
		Parameters:    file.FileCreationDefinition.Parameter,
		CommandFn:     file.FileCreationDefinition.Function,
	}
}
