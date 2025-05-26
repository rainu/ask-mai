package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/command"
	"log/slog"
	"mvdan.cc/sh/v3/shell"
	"strings"
)

type Command string

func (c Command) Validate() error {
	if len(c) == 0 {
		return fmt.Errorf("empty command")
	}

	return nil
}

func (c Command) CommandFn(fd FunctionDefinition) CommandFn {
	return func(ctx context.Context, argsAsJson string) ([]byte, error) {
		cmdDesc := command.CommandDescriptor{}
		var err error

		cmdDesc.Command, cmdDesc.Arguments, err = fd.GetCommandWithArgs(argsAsJson)
		if err != nil {
			return nil, fmt.Errorf("error creating command for tool '%s': %w", fd.Name, err)
		}

		if len(fd.Environment) > 0 {
			cmdDesc.Environment, err = fd.GetEnvironment(argsAsJson)
			if err != nil {
				return nil, fmt.Errorf("error creating environment for tool '%s': %w", fd.Name, err)
			}
		}
		if len(fd.AdditionalEnvironment) > 0 {
			cmdDesc.AdditionalEnvironment, err = fd.GetAdditionalEnvironment(argsAsJson)
			if err != nil {
				return nil, fmt.Errorf("error creating additional environment for tool '%s': %w", fd.Name, err)
			}
		}
		if fd.WorkingDir != "" {
			cmdDesc.WorkingDirectory, err = fd.GetWorkingDirectory(argsAsJson)
			if err != nil {
				return nil, fmt.Errorf("error creating working directory for tool '%s': %w", fd.Name, err)
			}
		}

		return cmdDesc.Run(ctx)
	}
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

func (f *FunctionDefinition) GetEnvironment(jsonArgs string) (map[string]string, error) {
	return processEnv(f.Environment, jsonArgs)
}

func (f *FunctionDefinition) GetAdditionalEnvironment(jsonArgs string) (map[string]string, error) {
	return processEnv(f.AdditionalEnvironment, jsonArgs)
}

func processEnv(env map[string]string, jsonArgs string) (map[string]string, error) {
	var data parsedArgs
	if err := json.Unmarshal([]byte(jsonArgs), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	result := map[string]string{}
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
