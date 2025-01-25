package config

import (
	"errors"
	"fmt"
	"github.com/rainu/ask-mai/config/expression"
	"github.com/wailsapp/wails/v2/pkg/options"
)

const (
	TranslucentNever = "never"
	TranslucentEver  = "ever"
	TranslucentHover = "hover"
)

type WindowConfig struct {
	Title            string                     `yaml:"title" usage:"The window title"`
	InitialWidth     expression.NumberContainer `yaml:"init-width" usage:"Expression: The (initial) width of the window"`
	MaxHeight        expression.NumberContainer `yaml:"max-height" usage:"Expression: The maximal height of the chat response area"`
	InitialPositionX expression.NumberContainer `yaml:"init-pos-x" usage:"Expression: The (initial) x-position of the window"`
	InitialPositionY expression.NumberContainer `yaml:"init-pos-y" usage:"Expression: The (initial) y-position of the window (if grow-top is set, the y-position is inverted -> 0 is bottom instead of top)"`
	InitialZoom      expression.NumberContainer `yaml:"init-zoom" usage:"Expression: The (initial) zoom level of the window"`
	BackgroundColor  WindowBackgroundColor      `yaml:"bg-color"`
	StartState       int                        `yaml:"start-state"`
	AlwaysOnTop      bool                       `yaml:"always-on-top" usage:"Should the window be always on top"`
	GrowTop          bool                       `yaml:"grow-top" usage:"Should the window grow from bottom to the top"`
	Frameless        bool                       `yaml:"frameless" usage:"Should the window be frameless"`
	Resizeable       bool                       `yaml:"resizeable" usage:"Should the window be resizeable"`
	Translucent      string                     `yaml:"translucent"`
}

type WindowBackgroundColor struct {
	R uint `yaml:"r" usage:"red value"`
	G uint `yaml:"g" usage:"green value"`
	B uint `yaml:"b" usage:"blue value"`
	A uint `yaml:"a" usage:"alpha value"`
}

func (w *WindowConfig) GetUsage(field string) string {
	switch field {
	case "StartState":
		return fmt.Sprintf("The window start state (normal(%d), minimized(%d), maximized(%d), fullscreen(%d))", options.Normal, options.Minimised, options.Maximised, options.Fullscreen)
	case "Translucent":
		return fmt.Sprintf("When the window should be translucent (%s, %s, %s)", TranslucentNever, TranslucentEver, TranslucentHover)
	}
	return ""
}

func (w *WindowConfig) Validate() error {

	if w.BackgroundColor.R > 255 {
		return fmt.Errorf("Invalid background color (red)")
	}
	if w.BackgroundColor.G > 255 {
		return fmt.Errorf("Invalid background color (green)")
	}
	if w.BackgroundColor.B > 255 {
		return fmt.Errorf("Invalid background color (blue)")
	}
	if w.BackgroundColor.A > 255 {
		return fmt.Errorf("Invalid background color (alpha)")
	}
	if w.StartState < int(options.Normal) || w.StartState > int(options.Fullscreen) {
		return fmt.Errorf("Invalid window start state")
	}

	if err := expression.NumberExpression(w.MaxHeight.Expression).Validate(); err != nil {
		return fmt.Errorf("Invalid window max height expression: %w", err)
	}

	if err := expression.NumberExpression(w.InitialWidth.Expression).Validate(); err != nil {
		return fmt.Errorf("Invalid window initial width expression: %w", err)
	}

	if err := expression.NumberExpression(w.InitialPositionX.Expression).Validate(); err != nil {
		return fmt.Errorf("Invalid window initial x-position expression: %w", err)
	}

	if err := expression.NumberExpression(w.InitialPositionY.Expression).Validate(); err != nil {
		return fmt.Errorf("Invalid window initial y-position expression: %w", err)
	}

	if err := expression.NumberExpression(w.InitialZoom.Expression).Validate(); err != nil {
		return fmt.Errorf("Invalid window initial zoom expression: %w", err)
	}

	if w.Translucent != TranslucentNever && w.Translucent != TranslucentEver && w.Translucent != TranslucentHover {
		return fmt.Errorf("Invalid window translucent value")
	}

	return nil
}

func (w *WindowConfig) ResolveExpressions(variables expression.Variables) (err error) {
	var curErr error

	w.InitialWidth.Value, curErr = expression.NumberExpression(w.InitialWidth.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial width expression: %w", curErr))
	}
	w.MaxHeight.Value, curErr = expression.NumberExpression(w.MaxHeight.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving max width expression: %w", curErr))
	}
	w.InitialPositionX.Value, curErr = expression.NumberExpression(w.InitialPositionX.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial x-position expression: %w", curErr))
	}
	w.InitialZoom.Value, curErr = expression.NumberExpression(w.InitialZoom.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial zoom expression: %w", curErr))
	}

	w.InitialPositionY.Value, curErr = expression.NumberExpression(w.InitialPositionY.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial y-position expression: %w", curErr))
	}

	// in case of grow top the axis is inverted (0,0 is left bottom - instead of left top)
	if w.GrowTop {
		w.InitialPositionY.Value = float64(variables.CurrentScreen.Dimension.Height) - w.InitialPositionY.Value
	}

	return
}
