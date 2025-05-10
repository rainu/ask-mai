package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

const FunctionArgumentNameAll = "@"

type Config struct {
	Tools map[string]FunctionDefinition `yaml:"functions,omitempty" usage:"Function definition: "`

	BuiltInTools BuiltIns `yaml:"builtin,omitempty" usage:"Built-in tools: "`
}

type CommandFn func(ctx context.Context, jsonArguments string) ([]byte, error)
type ApprovalFn func(ctx context.Context, jsonArguments string) bool

type FunctionDefinition struct {
	Name        string `yaml:"-" json:"name" usage:"The name of the function"`
	Description string `yaml:"description,omitempty" json:"description" usage:"The description of the function"`
	Parameters  any    `yaml:"parameters,omitempty" json:"parameters" usage:"The parameter definition of the function"`
	Approval    string `yaml:"approval,omitempty" json:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	Command               string            `yaml:"command,omitempty,omitempty" json:"command,omitempty" usage:"The command to execute. This is a format string with placeholders for the parameters. Example: /usr/bin/touch $path"`
	CommandExpr           string            `yaml:"commandExpr,omitempty,omitempty" json:"commandExpr,omitempty" usage:"JavaScript expression (or path to JS-file) to execute. See Tool-Help (--help-tool) for more information."`
	Environment           map[string]string `yaml:"env,omitempty,omitempty" json:"env,omitempty" usage:"Environment variables to pass to the command (will overwrite the default environment)"`
	AdditionalEnvironment map[string]string `yaml:"additionalEnv,omitempty,omitempty" json:"additionalEnv,omitempty" usage:"Additional environment variables to pass to the command (will be merged with the default environment)"`
	WorkingDir            string            `yaml:"workingDir,omitempty,omitempty" json:"workingDir,omitempty" usage:"The working directory for the command"`

	//will be filled at runtime (and should not be filled by user in any way)
	isBuiltIn  bool
	CommandFn  CommandFn  `yaml:"-" json:"-"`
	ApprovalFn ApprovalFn `yaml:"-" json:"-"`
}

func (t *Config) SetDefaults() {
	t.BuiltInTools = *NewBuiltIns()
}

func (t *Config) Validate() error {
	for cmd, definition := range t.Tools {
		if definition.CommandExpr != "" {
			if ve := CommandExpression(definition.CommandExpr).Validate(); ve != nil {
				return ve
			}
			definition.CommandFn = CommandExpression(definition.CommandExpr).CommandFn(definition)
		} else if definition.Command != "" {
			if ve := Command(definition.Command).Validate(); ve != nil {
				return ve
			}
			definition.CommandFn = Command(definition.Command).CommandFn(definition)
		} else {
			return fmt.Errorf("Command for tool '%s' is missing", cmd)
		}
	}

	return nil
}

func (t *Config) GetTools() map[string]FunctionDefinition {
	allFunctions := map[string]FunctionDefinition{}

	for _, fd := range t.BuiltInTools.AsFunctionDefinitions() {
		fd.isBuiltIn = true
		allFunctions[fd.Name] = fd
	}

	for name, tool := range t.Tools {
		allFunctions[name] = tool
	}

	return allFunctions
}

type parsedArgs map[string]interface{}

func (p parsedArgs) Get(varName string) (string, error) {
	varValue, exists := p[varName]
	if !exists {
		return "", nil
	}

	val, err := json.Marshal(varValue)
	if err != nil {
		return "", err
	}
	sVal := string(val)
	if len(sVal) > 0 && sVal[0] == '"' {
		sVal = sVal[1:]
	}
	if len(sVal) > 0 && sVal[len(sVal)-1] == '"' {
		sVal = sVal[:len(sVal)-1]
	}
	return sVal, nil
}

func (f *FunctionDefinition) NeedApproval(ctx context.Context, jsonArgs string) bool {
	if f.ApprovalFn == nil {
		return Approval(f.Approval).NeedsApproval(ctx, jsonArgs, f)
	}
	return f.ApprovalFn(ctx, jsonArgs)
}

func (f *FunctionDefinition) IsBuiltIn() bool {
	return f.isBuiltIn
}
