package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	cmdchain "github.com/rainu/go-command-chain"
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
		cmdDesc := CommandDescriptor{}
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

type CommandDescriptor struct {
	Command               string            `json:"command"`
	Arguments             []string          `json:"arguments"`
	Environment           map[string]string `json:"env"`
	AdditionalEnvironment map[string]string `json:"additionalEnv"`
	WorkingDirectory      string            `json:"workingDir"`
}

func (c CommandDescriptor) Run(ctx context.Context) ([]byte, error) {
	cmdBuild := cmdchain.Builder().JoinWithContext(ctx, c.Command, c.Arguments...)

	if len(c.Environment) > 0 {
		cmdBuild = cmdBuild.WithEnvironmentMap(toAnyMap(c.Environment))
	}
	if len(c.AdditionalEnvironment) > 0 {
		cmdBuild = cmdBuild.WithAdditionalEnvironmentMap(toAnyMap(c.AdditionalEnvironment))
	}
	if c.WorkingDirectory != "" {
		cmdBuild = cmdBuild.WithWorkingDirectory(c.WorkingDirectory)
	}

	buf := bytes.NewBuffer([]byte{})
	execErr := cmdBuild.Finalize().
		WithOutput(buf).
		WithError(buf).
		Run()

	return buf.Bytes(), execErr
}

func toAnyMap(m map[string]string) map[any]any {
	result := map[any]any{}
	for k, v := range m {
		result[k] = v
	}
	return result
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
