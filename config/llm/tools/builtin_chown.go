package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type ChangeOwner struct {
	Disable    bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval bool `yaml:"no-approval" json:"no-approval" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Y ChangeOwnerResult    `config:"-"`
	Z ChangeOwnerArguments `config:"-"`
}

func (f ChangeOwner) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "changeOwner",
		Description: "Changes the owner of file or directory on the user's system. Does not work on 'Windows' or 'Plan 9' operating systems.",
		CommandFn:   f.Command,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "The path to the file or directory to change the owner for. Use '~' as placeholder for the user's home directory.",
				},
				"user_id": map[string]any{
					"type":        "number",
					"description": "The id of the user which should own the file or directory. Use -1 to keep the current user.",
				},
				"group_id": map[string]any{
					"type":        "number",
					"description": "The id of the group which should own the file or directory. Use -1 to keep the current group.",
				},
			},
			"additionalProperties": false,
			"required":             []string{"path", "user_id", "group_id"},
		},
		NeedsApproval: !f.NoApproval,
	}
}

type ChangeOwnerArguments struct {
	Path Path `json:"path"`
	Uid  int  `json:"user_id"`
	Gid  int  `json:"group_id"`
}

type ChangeOwnerResult struct {
}

func (f ChangeOwner) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	var pArgs ChangeOwnerArguments
	err := json.Unmarshal([]byte(jsonArguments), &pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	if string(pArgs.Path) == "" {
		return nil, fmt.Errorf("missing parameter: 'path'")
	}
	path, err := pArgs.Path.Get()
	if err != nil {
		return nil, err
	}

	err = os.Chown(path, pArgs.Uid, pArgs.Gid)
	if err != nil {
		return nil, fmt.Errorf("error changing owner: %w", err)
	}

	return json.Marshal(ChangeOwnerResult{})
}
