package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type Stats struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.StatsResult    `config:"-" yaml:"-"`
	Z file.StatsArguments `config:"-" yaml:"-"`
}

func NewStats() Stats {
	return Stats{
		Approval: ApprovalNever,
	}
}

func (f Stats) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getStats",
		Approval:    f.Approval,
		Description: file.StatsDefinition.Description,
		Parameters:  file.StatsDefinition.Parameter,
		CommandFn:   file.StatsDefinition.Function,
	}
}
