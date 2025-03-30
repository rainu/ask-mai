package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/rainu/ask-mai/config/expression"
	"github.com/rainu/ask-mai/llms/tools/command"
	"log/slog"
	"strings"
)

type CommandExecution struct {
	Disable                bool     `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval             bool     `config:"no-approval" yaml:"no-approval" usage:"Needs no user approval to be executed"`
	NoApprovalCommands     []string `config:"allow" yaml:"allow" usage:"Needs no user approval for the specific command to be executed"`
	NoApprovalCommandsExpr []string `config:"allow-expr" yaml:"allow-expr" usage:"Needs no user approval for the specific command-line to be executed\nJavaScript expression - Return true if the command should be allowed:\n\tv.Name : string - contains the commands name\n\tv.Arguments : []string - contains the arguments\n\tv.WorkingDirectory: string - contains the working directory\n\tv.Environment : map[string]string - contains the environment variables\nExamples:\n\tv.Name == 'ls' && v.Arguments.length == 0\n\tv.Name == 'find' && v.Arguments.findIndex(a => a === \"-exec\") == -1"`

	//only for wails to generate TypeScript types
	Z command.CommandExecutionArguments `config:"-" yaml:"-"`
}

func (c CommandExecution) AsFunctionDefinition() *FunctionDefinition {
	if c.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "executeCommand",
		ApprovalFn:  c.CheckApproval,
		Description: command.CommandExecutionDefinition.Description,
		Parameters:  command.CommandExecutionDefinition.Parameter,
		CommandFn:   command.CommandExecutionDefinition.Function,
	}
}

func (c CommandExecution) CheckApproval(ctx context.Context, jsonArguments string) bool {
	// no command needs an approval
	if c.NoApproval {
		return false
	}

	var pArgs command.CommandExecutionArguments
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

func CalcApprovalExpr(e string, v command.CommandExecutionArguments) (bool, error) {
	vm := goja.New()
	err := vm.Set(expression.VarNameVariables, v)
	if err != nil {
		return false, fmt.Errorf("error setting variables: %w", err)
	}
	err = expression.SetupLog(vm)
	if err != nil {
		return false, fmt.Errorf("error setting functions: %w", err)
	}
	result, err := vm.RunString(e)
	if err != nil {
		return false, fmt.Errorf("error running expression: %w", err)
	}

	return result.ToBoolean(), nil
}
