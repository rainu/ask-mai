package tools

import (
	"context"
	"time"
)

type SystemTime struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `yaml:"approval" json:"approval" usage:"Needs user approval to be executed"`
}

func (s SystemTime) AsFunctionDefinition() *FunctionDefinition {
	if s.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getSystemTime",
		Description: "Get the current system time.",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
			"required":   []string{},
		},
		CommandFn:     s.Command,
		NeedsApproval: s.NeedsApproval,
	}
}

func (s SystemTime) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	return []byte(time.Now().String()), nil
}
