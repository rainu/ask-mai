package tools

import (
	"github.com/rainu/ask-mai/internal/mcp/server/tools/system"
)

type SystemInfo struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y system.SystemInfoResult    `yaml:"-"`
	Z system.SystemInfoArguments `yaml:"-"`
}

func (c *SystemInfo) SetDefaults() {
	if c.Approval == "" {
		c.Approval = ApprovalNever
	}
}

func NewSystemInfo() SystemInfo {
	return SystemInfo{
		Approval: ApprovalNever,
	}
}

func (s SystemInfo) AsFunctionDefinition() *FunctionDefinition {
	if s.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getSystemInformation",
		Approval:    s.Approval,
		Description: system.SystemInfoDefinition.Description,
		Parameters:  system.SystemInfoDefinition.Parameter,
		CommandFn:   system.SystemInfoDefinition.Function,
	}
}
