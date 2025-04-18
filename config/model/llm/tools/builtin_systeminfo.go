package tools

import (
	"github.com/rainu/ask-mai/llms/tools"
)

type SystemInfo struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y tools.SystemInfoResult    `config:"-" yaml:"-"`
	Z tools.SystemInfoArguments `config:"-" yaml:"-"`
}

func (s SystemInfo) AsFunctionDefinition() *FunctionDefinition {
	if s.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "getSystemInformation",
		NeedsApproval: s.NeedsApproval,
		Description:   tools.SystemInfoDefinition.Description,
		Parameters:    tools.SystemInfoDefinition.Parameter,
		CommandFn:     tools.SystemInfoDefinition.Function,
	}
}
