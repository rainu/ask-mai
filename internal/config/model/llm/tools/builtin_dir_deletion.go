package tools

import (
	"github.com/rainu/ask-mai/internal/llms/tools/file"
)

type DirectoryDeletion struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.DirectoryDeletionResult    `yaml:"-"`
	Z file.DirectoryDeletionArguments `yaml:"-"`
}

func NewDirectoryDeletion() DirectoryDeletion {
	return DirectoryDeletion{
		Approval: ApprovalAlways,
	}
}

func (f DirectoryDeletion) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "deleteDirectory",
		Approval:    f.Approval,
		Description: file.DirectoryDeletionDefinition.Description,
		Parameters:  file.DirectoryDeletionDefinition.Parameter,
		CommandFn:   file.DirectoryDeletionDefinition.Function,
	}
}
