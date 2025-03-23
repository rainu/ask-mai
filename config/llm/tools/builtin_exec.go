package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/rainu/ask-mai/config/expression"
	cmdchain "github.com/rainu/go-command-chain"
	"log/slog"
	"strings"
)

type CommandExecution struct {
	Disable                bool     `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval             bool     `config:"no-approval" yaml:"no-approval" usage:"Needs no user approval to be executed"`
	NoApprovalCommands     []string `config:"allow" yaml:"allow" usage:"Needs no user approval for the specific command to be executed"`
	NoApprovalCommandsExpr []string `config:"allow-expr" yaml:"allow-expr" usage:"Needs no user approval for the specific command-line to be executed\nJavaScript expression - Return true if the command should be allowed:\n\tv.Name : string - contains the commands name\n\tv.Arguments : []string - contains the arguments\n\tv.WorkingDirectory: string - contains the working directory\n\tv.Environment : map[string]string - contains the environment variables\nExamples:\n\tv.Name == 'ls' && v.Arguments.length == 0\n\tv.Name == 'find' && v.Arguments.findIndex(a => a === \"-exec\") == -1"`

	//only for wails to generate TypeScript types
	Z CommandExecutionArguments `config:"-"`
}

func (c CommandExecution) AsFunctionDefinition() *FunctionDefinition {
	if c.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "executeCommand",
		Description: "Execute a command on the user's system.",
		Parameters: map[string]any{
			"type":   "object",
			"strict": true,
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "The name / path to the command to execute.",
				},
				"arguments": map[string]any{
					"type":        "array",
					"description": "The arguments for the command.",
					"items": map[string]any{
						"type": "string",
					},
				},
				"working_directory": map[string]any{
					"type":        "string",
					"description": "The working directory for the command.",
				},
				"environment": map[string]any{
					"type":                 "object",
					"description":          "Additional environment variables to pass to the command.",
					"additionalProperties": true,
				},
			},
			"required": []string{"name"},
		},
		CommandFn:  c.Command,
		ApprovalFn: c.CheckApproval,
	}
}

type CommandExecutionArguments struct {
	Name             string            `json:"name"`
	Arguments        []string          `json:"arguments"`
	WorkingDirectory string            `json:"working_directory"`
	Environment      map[string]string `json:"environment"`
}

func (c CommandExecution) CheckApproval(ctx context.Context, jsonArguments string) bool {
	// no command needs an approval
	if c.NoApproval {
		return false
	}

	var pArgs CommandExecutionArguments
	err := json.Unmarshal([]byte(jsonArguments), &pArgs)
	if err != nil {
		slog.Error("Error parsing argument!", "error", err)
		return true
	}

	split := strings.Split(pArgs.Name, "/")
	command := split[len(split)-1]
	command = strings.ToLower(command)
	command = strings.TrimSpace(command)

	for _, allowCommand := range c.NoApprovalCommands {
		allowCommand = strings.TrimSpace(allowCommand)
		allowCommand = strings.ToLower(allowCommand)

		if command == allowCommand {
			return false
		}
	}

	cmdLine := strings.Join(append([]string{pArgs.Name}, pArgs.Arguments...), " ")
	for _, expr := range c.NoApprovalCommandsExpr {
		allowed, err := CalcApprovalExpr(expr, pArgs)
		if err != nil {
			slog.Error("Error calculating expression!", "error", err)
			continue
		}

		if allowed {
			slog.Debug("Command allowed by expression!",
				"command", cmdLine,
				"expression", expr,
			)
			return false
		}
	}

	// needs approval
	return true
}

func CalcApprovalExpr(e string, v CommandExecutionArguments) (bool, error) {
	vm := goja.New()
	err := vm.Set(expression.VarNameVariables, v)
	if err != nil {
		return false, fmt.Errorf("error setting variables: %w", err)
	}
	err = expression.SetupLog(vm)
	if err != nil {
		return false, fmt.Errorf("error setting functions: %w", err)
	}
	result, err := vm.RunString(string(e))
	if err != nil {
		return false, fmt.Errorf("error running expression: %w", err)
	}

	return result.ToBoolean(), nil
}

func (c CommandExecution) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	var pArgs CommandExecutionArguments
	err := json.Unmarshal([]byte(jsonArguments), &pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	if pArgs.Name == "" {
		return nil, fmt.Errorf("missing parameter: 'name'")
	}

	cmd := cmdchain.Builder().JoinWithContext(ctx, pArgs.Name, pArgs.Arguments...)

	if pArgs.WorkingDirectory != "" {
		cmd = cmd.WithWorkingDirectory(pArgs.WorkingDirectory)
	}
	if len(pArgs.Environment) > 0 {
		envMap := map[any]any{}
		for k, v := range pArgs.Environment {
			envMap[k] = v
		}
		cmd = cmd.WithAdditionalEnvironmentMap(envMap)
	}

	buf := bytes.NewBuffer([]byte{})
	execErr := cmd.Finalize().
		WithOutput(buf).
		WithError(buf).
		Run()

	return buf.Bytes(), execErr
}
