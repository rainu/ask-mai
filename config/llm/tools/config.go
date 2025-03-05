package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"log/slog"
	"mvdan.cc/sh/v3/shell"
	"strings"
)

const FunctionArgumentNameAll = "@"

type Config struct {
	RawTools []string `config:"function" yaml:"-" usage:"Function definition (json) to use. See Tool-Help (--help-tool) for more information."`

	Tools map[string]FunctionDefinition `config:"-" yaml:"functions"`

	BuiltInTools BuiltIns `config:"builtin" yaml:"builtin" usage:"Built-in tools: "`
}

type CommandFn func(ctx context.Context, jsonArguments string) ([]byte, error)

type FunctionDefinition struct {
	Name          string `config:"name" yaml:"-" json:"name" usage:"The name of the function"`
	Description   string `yaml:"description" json:"description" usage:"The description of the function"`
	Parameters    any    `yaml:"parameters" json:"parameters" usage:"The parameter definition of the function"`
	NeedsApproval bool   `yaml:"approval" json:"approval" usage:"Needs user approval to be executed"`

	Command               string            `yaml:"command" json:"command" usage:"The command to execute. This is a format string with placeholders for the parameters. Example: /usr/bin/touch $path"`
	Environment           map[string]string `yaml:"env,omitempty" json:"env,omitempty" usage:"Environment variables to pass to the command (will overwrite the default environment)"`
	AdditionalEnvironment map[string]string `yaml:"additionalEnv,omitempty" json:"additionalEnv,omitempty" usage:"Additional environment variables to pass to the command (will be merged with the default environment)"`
	WorkingDir            string            `yaml:"workingDir,omitempty" json:"workingDir,omitempty" usage:"The working directory for the command"`

	CommandFn CommandFn `config:"-" yaml:"-" json:"-"` // BuiltIn functions
}

func (t *Config) Validate() error {
	for i, tool := range t.RawTools {
		var result FunctionDefinition

		err := json.Unmarshal([]byte(tool), &result)
		if err != nil {
			return fmt.Errorf("Invalid tool definition #%d: %w", i, err)
		}

		t.Tools[result.Name] = result
	}

	for cmd, definition := range t.Tools {
		if definition.Command == "" {
			return fmt.Errorf("Command for tool '%s' is missing", cmd)
		}
	}

	return nil
}

func (t *Config) AsOptions() (opts []llms.CallOption) {
	var tools []llms.Tool

	for name, definition := range t.GetTools() {
		tools = append(tools, llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        name,
				Description: definition.Description,
				Parameters:  definition.Parameters,
			},
		})
	}

	if len(tools) > 0 {
		opts = append(opts, llms.WithTools(tools))
	}
	return
}

func (t *Config) GetTools() map[string]FunctionDefinition {
	allFunctions := map[string]FunctionDefinition{}

	for _, fd := range t.BuiltInTools.AsFunctionDefinitions() {
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

func (f *FunctionDefinition) IsBuiltIn() bool {
	return f.CommandFn != nil
}

func (f *FunctionDefinition) GetCommandWithArgs(jsonArgs string) (string, []string, error) {
	var data parsedArgs
	if err := json.Unmarshal([]byte(jsonArgs), &data); err != nil {
		return "", nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	fields, err := shell.Fields(f.Command, func(varName string) string {
		if varName == FunctionArgumentNameAll {
			return jsonArgs
		}
		r, err := data.Get(varName)
		if err != nil {
			slog.Error("Failed to marshal value",
				"varName", varName,
				"value", r,
				"error", err,
			)
		}
		return r
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse command: %w", err)
	}
	return fields[0], fields[1:], nil
}

func (f *FunctionDefinition) GetEnvironment(jsonArgs string) (map[any]any, error) {
	return processEnv(f.Environment, jsonArgs)
}

func (f *FunctionDefinition) GetAdditionalEnvironment(jsonArgs string) (map[any]any, error) {
	return processEnv(f.AdditionalEnvironment, jsonArgs)
}

func processEnv(env map[string]string, jsonArgs string) (map[any]any, error) {
	var data parsedArgs
	if err := json.Unmarshal([]byte(jsonArgs), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	result := map[any]any{}
	for key, value := range env {
		for vk := range data {
			v, err := data.Get(vk)
			if err != nil {
				return nil, err
			}
			value = strings.Replace(value, "$"+vk, v, -1)
		}
		value = strings.Replace(value, "$@", jsonArgs, -1)
		result[key] = value
	}

	return result, nil
}

func (f *FunctionDefinition) GetWorkingDirectory(jsonArgs string) (string, error) {
	var data parsedArgs
	if err := json.Unmarshal([]byte(jsonArgs), &data); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	value := f.WorkingDir
	for vk := range data {
		v, err := data.Get(vk)
		if err != nil {
			return "", err
		}
		value = strings.Replace(value, "$"+vk, v, -1)
	}

	return value, nil
}
