package expression

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
)

func initRuntime(ctx context.Context) (vm *goja.Runtime, err error) {
	vm = goja.New()

	//setup functions
	err = vm.Set(FuncNameLog, Log)
	if err != nil {
		return nil, fmt.Errorf("unable to set %s function: %w", FuncNameLog, err)
	}

	err = vm.Set(FuncNameRun, run(ctx, vm))
	if err != nil {
		return nil, fmt.Errorf("unable to set %s function: %w", FuncNameRun, err)
	}

	err = vm.Set(FuncNameFetch, fetch(ctx, vm))
	if err != nil {
		return nil, fmt.Errorf("unable to set %s function: %w", FuncNameFetch, err)
	}

	//setup global variables
	for key, value := range globalVariables {
		err = vm.Set(key, value)
		if err != nil {
			return nil, fmt.Errorf("unable to set %s variable: %w", key, err)
		}
	}

	return
}

func Run(ctx context.Context, expression string, ctxVal any) *Result {
	vm, err := initRuntime(ctx)
	if err != nil {
		return &Result{err: fmt.Errorf("unable to initialize runtime: %w", err)}
	}

	// set additional variable
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	err = vm.Set(VarNameContext, ctxVal)
	if err != nil {
		return &Result{err: fmt.Errorf("unable to set %s variable: %w", VarNameContext, err)}
	}

	var v goja.Value

	prog, wasPrecompiled := precompiledPrograms[expression]
	if wasPrecompiled {
		v, err = vm.RunProgram(prog)
	} else {
		v, err = vm.RunString(expression)
	}

	if err != nil {
		return &Result{err: fmt.Errorf("unable to run expression: %w", err)}
	}

	return &Result{result: v}
}

func Validate(expression string) error {
	_, err := goja.Compile("", expression, false)
	if err != nil {
		return fmt.Errorf("unable to compile expression: %w", err)
	}

	return nil
}
