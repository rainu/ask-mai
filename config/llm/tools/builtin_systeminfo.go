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

	//only for wails to generate TypeScript types
	Y SystemInfoResult    `config:"-"`
	Z SystemInfoArguments `config:"-"`
}

func (s SystemInfo) AsFunctionDefinition() *FunctionDefinition {
	if s.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getSystemInformation",
		Description: "Get the following information about the user's system: OS, architecture, number of CPUs, hostname, user directory, user ID, group ID, working directory, process ID.",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
			"required":   []string{},
		},
		CommandFn:     s.Command,
		NeedsApproval: s.NeedsApproval,
	}
}

type SystemInfoArguments struct {
}

type SystemInfoResult struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	CPU      int    `json:"cpus"`
	Hostname string `json:"hostname"`
	UserDir  string `json:"user_dir"`
	UserId   int    `json:"user_id"`
	GroupId  int    `json:"group_id"`
	WorkDir  string `json:"working_directory"`
	PID      int    `json:"process_id"`
}

func (s SystemInfo) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	return json.Marshal(SystemInfoResult{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		CPU:  runtime.NumCPU(),
		Hostname: func() string {
			h, err := os.Hostname()
			if err != nil {
				return "unknown"
			}
			return h
		}(),
		UserDir: func() string {
			home, err := os.UserHomeDir()
			if err != nil {
				return "unknown"
			}
			return home
		}(),
		UserId:  os.Getuid(),
		GroupId: os.Getgid(),
		WorkDir: func() string {
			dir, err := os.Getwd()
			if err != nil {
				return "unknown"
			}
			return dir
		}(),
		PID: os.Getpid(),
	})
}
