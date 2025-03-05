package tools

import (
	"context"
	"encoding/json"
	"os"
	"runtime"
)

type SystemInfo struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `yaml:"approval" json:"approval" usage:"Needs user approval to be executed"`
}

func (s SystemInfo) AsFunctionDefinition() *FunctionDefinition {
	if s.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "getSystemInformation",
		Description:   "Get some information about the user's system.",
		CommandFn:     s.Command,
		NeedsApproval: s.NeedsApproval,
	}
}

func (s SystemInfo) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	info := map[string]interface{}{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
		"cpus": runtime.NumCPU(),
		"hostname": func() string {
			h, err := os.Hostname()
			if err != nil {
				return "unknown"
			}
			return h
		}(),
		"user_dir": func() string {
			home, err := os.UserHomeDir()
			if err != nil {
				return "unknown"
			}
			return home
		}(),
		"user_id":  os.Getuid(),
		"group_id": os.Getgid(),
		"working_directory": func() string {
			dir, err := os.Getwd()
			if err != nil {
				return "unknown"
			}
			return dir
		}(),
		"process_id": os.Getpid(),
	}

	return json.Marshal(info)
}
