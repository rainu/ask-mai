package config

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/options"
)

const (
	TranslucentNever = "never"
	TranslucentEver  = "ever"
	TranslucentHover = "hover"
)

type WindowConfig struct {
	Title            string                `config:"title" usage:"The window title"`
	InitialWidth     ExpressionContainer   `config:"init-width" usage:"Expression: The (initial) width of the window"`
	MaxHeight        ExpressionContainer   `config:"max-height" usage:"Expression: The maximal height of the chat response area"`
	InitialPositionX ExpressionContainer   `config:"init-pos-x" usage:"Expression: The (initial) x-position of the window"`
	InitialPositionY ExpressionContainer   `config:"init-pos-y" usage:"Expression: The (initial) y-position of the window"`
	InitialZoom      ExpressionContainer   `config:"init-zoom" usage:"Expression: The (initial) zoom level of the window"`
	BackgroundColor  WindowBackgroundColor `config:"bg-color"`
	StartState       int                   `config:"start-state"`
	AlwaysOnTop      bool                  `config:"always-on-top" usage:"Should the window be always on top"`
	Frameless        bool                  `config:"frameless" usage:"Should the window be frameless"`
	Resizeable       bool                  `config:"resizeable" usage:"Should the window be resizeable"`
	Translucent      string                `config:"translucent"`
}

type WindowBackgroundColor struct {
	R uint `config:"r" usage:"red value"`
	G uint `config:"g" usage:"green value"`
	B uint `config:"b" usage:"blue value"`
	A uint `config:"a" usage:"alpha value"`
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

	if w.MaxHeight.Expression != "" {
		if err := ValidateExpression(w.MaxHeight.Expression); err != nil {
			return fmt.Errorf("Invalid window max height expression: %w", err)
		}
	}

	if w.InitialWidth.Expression != "" {
		if err := ValidateExpression(w.InitialWidth.Expression); err != nil {
			return fmt.Errorf("Invalid window initial width expression: %w", err)
		}
	}

	if w.InitialPositionX.Expression != "" {
		if err := ValidateExpression(w.InitialPositionX.Expression); err != nil {
			return fmt.Errorf("Invalid window initial x-position expression: %w", err)
		}
	}

	if w.InitialPositionY.Expression != "" {
		if err := ValidateExpression(w.InitialPositionY.Expression); err != nil {
			return fmt.Errorf("Invalid window initial y-position expression: %w", err)
		}
	}

	if w.InitialZoom.Expression != "" {
		if err := ValidateExpression(w.InitialZoom.Expression); err != nil {
			return fmt.Errorf("Invalid window initial zoom expression: %w", err)
		}
	}

	if w.Translucent != TranslucentNever && w.Translucent != TranslucentEver && w.Translucent != TranslucentHover {
		return fmt.Errorf("Invalid window translucent value")
	}

	return nil
}
