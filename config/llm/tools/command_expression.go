package tools

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/rainu/ask-mai/config/expression"
	"github.com/rainu/ask-mai/llms/tools/command"
	"os"
)

const FuncNameRun = "runCommand"

type CommandExpression string

type CommandVariables struct {
	FunctionDefinition FunctionDefinition `json:"fd"`
	Arguments          string             `json:"args"`
}

var parsedPrograms = map[string]*goja.Program{}

func (c CommandExpression) Validate() error {
	if len(c) == 0 {
		return nil
	}

	file, err := os.Open(string(c))
	if err == nil && file != nil {
		defer file.Close()

		ast, err := parser.ParseFile(nil, file.Name(), file, 0)
		if err != nil {
			return fmt.Errorf("error parsing file: %w", err)
		}
		prog, err := goja.CompileAST(ast, false)
		if err != nil {
			return fmt.Errorf("error compiling file: %w", err)
		}
		parsedPrograms[string(c)] = prog
	} else {
		prog, err := goja.Compile("", string(c), false)
		if err != nil {
			return fmt.Errorf("error compiling file: %w", err)
		}
		parsedPrograms[string(c)] = prog
	}

	return nil
}

func (c CommandExpression) CommandFn(fd FunctionDefinition) CommandFn {
	return func(ctx context.Context, args string) ([]byte, error) {
		vm := goja.New()
		vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

		err := vm.Set(expression.VarNameVariables, CommandVariables{
			FunctionDefinition: fd,
			Arguments:          args,
		})
		if err != nil {
			return nil, fmt.Errorf("error setting variables: %w", err)
		}
		err = expression.SetupLog(vm)
		if err != nil {
			return nil, fmt.Errorf("error setting functions: %w", err)
		}
		err = vm.Set(FuncNameRun, runCommand(ctx, vm))
		if err != nil {
			return nil, fmt.Errorf("error setting functions: %w", err)
		}

		prog := parsedPrograms[string(c)]

		result, err := vm.RunProgram(prog)
		if err != nil {
			return nil, fmt.Errorf("error running expression: %w", err)
		}

		return []byte(result.String()), nil
	}
}

func runCommand(ctx context.Context, vm *goja.Runtime) func(command.CommandDescriptor) string {
	return func(cmd command.CommandDescriptor) string {
		r, err := cmd.Run(ctx)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		return string(r)
	}
}
