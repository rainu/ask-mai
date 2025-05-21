package tools

import (
	"github.com/rainu/ask-mai/internal/mcp/server/tools/file"
)

type ChangeOwner struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.ChangeOwnerResult    `yaml:"-"`
	Z file.ChangeOwnerArguments `yaml:"-"`
}

func (c *ChangeOwner) SetDefaults() {
	if c.Approval == "" {
		c.Approval = ApprovalAlways
	}
}

func NewChangeOwner() ChangeOwner {
	return ChangeOwner{
		Approval: ApprovalAlways,
	}
}

func (f ChangeOwner) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "changeOwner",
		Approval:    f.Approval,
		Description: file.ChangeOwnerDefinition.Description,
		Parameters:  file.ChangeOwnerDefinition.Parameter,
		CommandFn:   file.ChangeOwnerDefinition.Function,
	}
}
