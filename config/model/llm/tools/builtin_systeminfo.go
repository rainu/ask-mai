package tools

import (
	"github.com/rainu/ask-mai/llms/tools"
)

type SystemInfo struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y tools.SystemInfoResult    `yaml:"-"`
	Z tools.SystemInfoArguments `yaml:"-"`
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
		Description: tools.SystemInfoDefinition.Description,
		Parameters:  tools.SystemInfoDefinition.Parameter,
		CommandFn:   tools.SystemInfoDefinition.Function,
	}
}
