package expression

import (
	"fmt"
	"github.com/dop251/goja"
)

type StringContainer struct {
	Expression string `config:"" yaml:"expression"`
	Value      string `config:"-" yaml:"value"`
}

type StringExpression string

func (e StringExpression) Calculate() (string, error) {
	vm := goja.New()
	err := setupLog(vm)
	if err != nil {
		return "", fmt.Errorf("error setting functions: %w", err)
	}
	result, err := vm.RunString(string(e))
	if err != nil {
		return "", fmt.Errorf("error running expression: %w", err)
	}

	return result.String(), nil
}

func (e StringExpression) Validate() error {
	_, err := e.Calculate()
	return err
}
