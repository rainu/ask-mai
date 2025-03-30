package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type DirectoryCreation struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.DirectoryCreationResult    `config:"-" yaml:"-"`
	Z file.DirectoryCreationArguments `config:"-" yaml:"-"`
}

func (f DirectoryCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "createDirectory",
		NeedsApproval: f.NeedsApproval,
		Description:   file.DirectoryCreationDefinition.Description,
		Parameters:    file.DirectoryCreationDefinition.Parameter,
		CommandFn:     file.DirectoryCreationDefinition.Function,
	}
}
