package expression

import (
	"fmt"
	"github.com/dop251/goja"
)

type NumberContainer struct {
	Expression string  `config:"" yaml:"expression"`
	Value      float64 `config:"-" yaml:"value"`
}

type NumberExpression string

func (e NumberExpression) Calculate(v Variables) (float64, error) {
	vm := goja.New()
	err := vm.Set(VarNameVariables, v)
	if err != nil {
		return 0, fmt.Errorf("error setting variables: %w", err)
	}
	err = SetupLog(vm)
	if err != nil {
		return 0, fmt.Errorf("error setting functions: %w", err)
	}
	result, err := vm.RunString(string(e))
	if err != nil {
		return 0, fmt.Errorf("error running expression: %w", err)
	}

	if result.ToNumber().SameAs(goja.NaN()) {
		return 0, fmt.Errorf("result is not a number")
	}

	return result.ToFloat(), nil
}

func (e NumberExpression) Validate() error {
	if e == "" {
		return nil
	}

	tv := Variables{}
	for i := 0; i < 1000; i++ {
		tv.Screens = append(tv.Screens, VariableScreen{
			Dimension: VariableScreenDimension{
				Width:  1920,
				Height: 1080,
			},
		})
	}
	tv.PrimaryScreen = tv.Screens[0]

	_, err := e.Calculate(tv)
	return err
}
