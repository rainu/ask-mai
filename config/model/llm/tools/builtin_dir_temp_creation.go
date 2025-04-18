package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type DirectoryTempCreation struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.DirectoryTempCreationResult    `config:"-" yaml:"-"`
	Z file.DirectoryTempCreationArguments `config:"-" yaml:"-"`
}

func (f DirectoryTempCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "createTempDirectory",
		NeedsApproval: f.NeedsApproval,
		Description:   file.DirectoryTempCreationDefinition.Description,
		Parameters:    file.DirectoryTempCreationDefinition.Parameter,
		CommandFn:     file.DirectoryTempCreationDefinition.Function,
	}
}
