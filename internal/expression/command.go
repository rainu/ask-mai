package expression

import (
	"context"
	"github.com/dop251/goja"
	"github.com/rainu/ask-mai/internal/llms/tools/command"
)

const FuncNameRun = "run"

func run(ctx context.Context, vm *goja.Runtime) func(command.CommandDescriptor) string {
	return func(cmd command.CommandDescriptor) string {
		r, err := cmd.Run(ctx)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		return string(r)
	}
}
