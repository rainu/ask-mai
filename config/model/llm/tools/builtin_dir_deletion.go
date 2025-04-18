package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type DirectoryDeletion struct {
	Disable    bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval bool `config:"no-approval" yaml:"no-approval" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.DirectoryDeletionResult    `config:"-" yaml:"-"`
	Z file.DirectoryDeletionArguments `config:"-" yaml:"-"`
}

func (f DirectoryDeletion) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "deleteDirectory",
		NeedsApproval: !f.NoApproval,
		Description:   file.DirectoryDeletionDefinition.Description,
		Parameters:    file.DirectoryDeletionDefinition.Parameter,
		CommandFn:     file.DirectoryDeletionDefinition.Function,
	}
}
