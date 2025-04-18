package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileDeletion struct {
	Disable    bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval bool `config:"no-approval" yaml:"no-approval" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.FileDeletionResult    `config:"-" yaml:"-"`
	Z file.FileDeletionArguments `config:"-" yaml:"-"`
}

func (f FileDeletion) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "deleteFile",
		NeedsApproval: !f.NoApproval,
		Description:   file.FileDeletionDefinition.Description,
		Parameters:    file.FileDeletionDefinition.Parameter,
		CommandFn:     file.FileDeletionDefinition.Function,
	}
}
