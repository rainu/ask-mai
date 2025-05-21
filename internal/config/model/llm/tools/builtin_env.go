package tools

import (
	"github.com/rainu/ask-mai/internal/mcp/server/tools/system"
)

type Environment struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y system.EnvironmentResult    `yaml:"-"`
	Z system.EnvironmentArguments `yaml:"-"`
}

func (c *Environment) SetDefaults() {
	if c.Approval == "" {
		c.Approval = ApprovalNever
	}
}

func NewEnvironment() Environment {
	return Environment{
		Approval: ApprovalNever,
	}
}

func (f Environment) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getEnvironment",
		Approval:    f.Approval,
		Description: system.EnvironmentDefinition.Description,
		Parameters:  system.EnvironmentDefinition.Parameter,
		CommandFn:   system.EnvironmentDefinition.Function,
	}
}
