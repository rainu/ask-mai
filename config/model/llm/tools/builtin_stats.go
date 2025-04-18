package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type Stats struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.StatsResult    `config:"-" yaml:"-"`
	Z file.StatsArguments `config:"-" yaml:"-"`
}

func (f Stats) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "getStats",
		NeedsApproval: f.NeedsApproval,
		Description:   file.StatsDefinition.Description,
		Parameters:    file.StatsDefinition.Parameter,
		CommandFn:     file.StatsDefinition.Function,
	}
}
