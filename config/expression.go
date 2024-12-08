package config

import (
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go/token"
	"go/types"
	"reflect"
	"strconv"
	"strings"
)

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

func (v Variables) ToFlatMap() (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	flatMap := make(map[string]interface{})
	flatten("", jsonData, flatMap)
	return flatMap, nil
}

func flatten(prefix string, data interface{}, flatMap map[string]interface{}) {
	if data == nil {
		return
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		for k, v := range data.(map[string]interface{}) {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			flatten(key, v, flatMap)
		}
	case reflect.Slice:
		for i, v := range data.([]interface{}) {
			key := fmt.Sprintf("%s[%d]", prefix, i)
			flatten(key, v, flatMap)
		}
	default:
		flatMap[prefix] = data
	}
}

func (e Expression) Calculate(v Variables) (float64, error) {
	flatVars, err := v.ToFlatMap()
	if err != nil {
		return 0, fmt.Errorf("error flattening variables: %w", err)
	}

	expression := string(e)
	for variable, value := range flatVars {
		expression = strings.ReplaceAll(expression, variable, fmt.Sprintf("%v", value))
	}

	fs := token.NewFileSet()
	tv, err := types.Eval(fs, nil, token.NoPos, expression)
	if err != nil {
		return 0, fmt.Errorf("error evaluating expression: %w", err)
	}

	result, err := strconv.ParseFloat(tv.Value.String(), 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing result: %w", err)
	}

	return result, nil
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
