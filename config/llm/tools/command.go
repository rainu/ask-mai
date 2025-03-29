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
	return nil
}

func (c Command) CommandFn(fd FunctionDefinition) CommandFn {
	return func(ctx context.Context, argsAsJson string) ([]byte, error) {
		cmd, args, err := fd.GetCommandWithArgs(argsAsJson)
		if err != nil {
			return nil, fmt.Errorf("error creating command for tool '%s': %w", fd.Name, err)
		}

		cmdBuild := cmdchain.Builder().JoinWithContext(ctx, cmd, args...)

		if len(fd.Environment) > 0 {
			env, err := fd.GetEnvironment(argsAsJson)
			if err != nil {
				return nil, fmt.Errorf("error creating environment for tool '%s': %w", fd.Name, err)
			}
			cmdBuild = cmdBuild.WithEnvironmentMap(env)
		}
		if len(fd.AdditionalEnvironment) > 0 {
			env, err := fd.GetAdditionalEnvironment(argsAsJson)
			if err != nil {
				return nil, fmt.Errorf("error creating additional environment for tool '%s': %w", fd.Name, err)
			}
			cmdBuild = cmdBuild.WithAdditionalEnvironmentMap(env)
		}
		if fd.WorkingDir != "" {
			wd, err := fd.GetWorkingDirectory(argsAsJson)
			if err != nil {
				return nil, fmt.Errorf("error creating working directory for tool '%s': %w", fd.Name, err)
			}
			cmdBuild = cmdBuild.WithWorkingDirectory(wd)
		}

		buf := bytes.NewBuffer([]byte{})
		execErr := cmdBuild.Finalize().
			WithOutput(buf).
			WithError(buf).
			Run()

		return buf.Bytes(), execErr
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
