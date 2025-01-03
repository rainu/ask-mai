package config

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

type ExpressionContainer struct {
	Expression string  `config:"" yaml:"expression"`
	Value      float64 `config:"-" yaml:"value"`
}

type Expression string

type Variables struct {
	PrimaryScreen    VariableScreen   `json:"PrimaryScreen"`
	CurrentScreen    VariableScreen   `json:"CurrentScreen"`
	Screens          []VariableScreen `json:"Screens"`
	SecondaryScreens []VariableScreen `json:"SecondaryScreens"`
}

type VariableScreen struct {
	Dimension VariableScreenDimension `json:"Dimension"`
}

type VariableScreenDimension struct {
	Width  int `json:"Width"`
	Height int `json:"Height"`
}

func FromScreens(screens []runtime.Screen) Variables {
	var variables Variables
	for _, screen := range screens {
		variables.Screens = append(variables.Screens, VariableScreen{
			Dimension: VariableScreenDimension{
				Width:  screen.PhysicalSize.Width,
				Height: screen.PhysicalSize.Height,
			},
		})
		if screen.IsPrimary {
			variables.PrimaryScreen = variables.Screens[len(variables.Screens)-1]
		} else {
			variables.SecondaryScreens = append(variables.SecondaryScreens, variables.Screens[len(variables.Screens)-1])
		}
		if screen.IsCurrent {
			variables.CurrentScreen = variables.Screens[len(variables.Screens)-1]
		}
	}

	return variables
}

const varNameVariables = "v"
const funcNameLog = "log"

func (e Expression) Calculate(v Variables) (float64, error) {
	vm := goja.New()
	err := vm.Set(varNameVariables, v)
	if err != nil {
		return 0, fmt.Errorf("error setting variables: %w", err)
	}
	err = vm.Set(funcNameLog, func(args ...interface{}) {
		fmt.Fprint(os.Stderr, "EXPRESSION_LOG: ")
		fmt.Fprintln(os.Stderr, args...)
	})
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

func ValidateExpression(e string) error {
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

	_, err := Expression(e).Calculate(tv)
	return err
}

func (c *Config) ResolveExpressions(variables Variables) (err error) {
	var curErr error

	c.UI.Window.InitialWidth.Value, curErr = Expression(c.UI.Window.InitialWidth.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial width expression: %w", curErr))
	}
	c.UI.Window.MaxHeight.Value, curErr = Expression(c.UI.Window.MaxHeight.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving max width expression: %w", curErr))
	}
	c.UI.Window.InitialPositionX.Value, curErr = Expression(c.UI.Window.InitialPositionX.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial x-position expression: %w", curErr))
	}
	c.UI.Window.InitialPositionY.Value, curErr = Expression(c.UI.Window.InitialPositionY.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial y-position expression: %w", curErr))
	}
	c.UI.Window.InitialZoom.Value, curErr = Expression(c.UI.Window.InitialZoom.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial zoom expression: %w", curErr))
	}

	return
}
