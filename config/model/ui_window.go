package model

import (
	"errors"
	"fmt"
	"github.com/rainu/ask-mai/config/model/common"
	"github.com/rainu/ask-mai/expression"
	"github.com/wailsapp/wails/v2/pkg/options"
	"runtime"
)

const (
	TranslucentNever = "never"
	TranslucentEver  = "ever"
	TranslucentHover = "hover"
)

type WindowConfig struct {
	Title            string                 `yaml:"title,omitempty" usage:"The window title"`
	InitialWidth     common.NumberContainer `yaml:"init-width,omitempty" usage:"Expression: The (initial) width of the window"`
	MaxHeight        common.NumberContainer `yaml:"max-height,omitempty" usage:"Expression: The maximal height of the chat response area"`
	InitialPositionX common.NumberContainer `yaml:"init-pos-x,omitempty" usage:"Expression: The (initial) x-position of the window"`
	InitialPositionY common.NumberContainer `yaml:"init-pos-y,omitempty" usage:"Expression: The (initial) y-position of the window (if grow-top is set, the y-position is inverted -> 0 is bottom instead of top)"`
	InitialZoom      common.NumberContainer `yaml:"init-zoom,omitempty" usage:"Expression: The (initial) zoom level of the window"`
	BackgroundColor  WindowBackgroundColor  `yaml:"bg-color,omitempty" usage:"The background color of the window: "`
	StartState       int                    `yaml:"start-state,omitempty"`
	AlwaysOnTop      bool                   `yaml:"always-on-top,omitempty" usage:"Should the window be always on top"`
	ShowTitleBar     bool                   `yaml:"show-title-bar,omitempty" usage:"Should the window show the title-bar"`
	TitleBarHeight   int                    `yaml:"title-bar-height,omitempty,omitempty" usage:"The height of the title bar"`
	GrowTop          bool                   `yaml:"grow-top,omitempty" usage:"Should the window grow from bottom to the top"`
	Frameless        bool                   `yaml:"frameless,omitempty" usage:"Should the window be frameless"`
	Resizeable       bool                   `yaml:"resizeable,omitempty" usage:"Should the window be resizeable"`
	Translucent      string                 `yaml:"translucent,omitempty"`
}

type WindowBackgroundColor struct {
	R uint `yaml:"r" usage:"red value"`
	G uint `yaml:"g" usage:"green value"`
	B uint `yaml:"b" usage:"blue value"`
	A uint `yaml:"a" usage:"alpha value"`
}

func (w *WindowConfig) SetDefaults() {
	w.Title = "Prompt - Ask mAI"
	w.InitialWidth = common.NumberContainer{Expression: expression.VarNameScreens + ".CurrentScreen.Dimension.Width/2"}
	w.MaxHeight = common.NumberContainer{Expression: expression.VarNameScreens + ".CurrentScreen.Dimension.Height/3"}
	w.InitialPositionX = common.NumberContainer{Expression: expression.VarNameScreens + ".CurrentScreen.Dimension.Width/4"}
	w.InitialPositionY = common.NumberContainer{Expression: "0"}
	w.InitialZoom = common.NumberContainer{Expression: "1.0"}
	w.BackgroundColor = WindowBackgroundColor{R: 255, G: 255, B: 255, A: 192}
	w.StartState = int(options.Normal)
	w.Translucent = TranslucentHover
	w.Frameless = true
	w.AlwaysOnTop = true
	w.Resizeable = true

	if runtime.GOOS == "windows" {
		w.TitleBarHeight = 32
	} else if runtime.GOOS == "darwin" {
		w.TitleBarHeight = 28
	}
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

	if err := expression.Validate(w.MaxHeight.Expression); err != nil {
		return fmt.Errorf("Invalid window max height expression: %w", err)
	}

	if err := expression.Validate(w.InitialWidth.Expression); err != nil {
		return fmt.Errorf("Invalid window initial width expression: %w", err)
	}

	if err := expression.Validate(w.InitialPositionX.Expression); err != nil {
		return fmt.Errorf("Invalid window initial x-position expression: %w", err)
	}

	if err := expression.Validate(w.InitialPositionY.Expression); err != nil {
		return fmt.Errorf("Invalid window initial y-position expression: %w", err)
	}

	if err := expression.Validate(w.InitialZoom.Expression); err != nil {
		return fmt.Errorf("Invalid window initial zoom expression: %w", err)
	}

	if w.Translucent != TranslucentNever && w.Translucent != TranslucentEver && w.Translucent != TranslucentHover {
		return fmt.Errorf("Invalid window translucent value")
	}

	return nil
}

func (w *WindowConfig) ResolveExpressions(screen expression.Screen) (err error) {
	var curErr error

	w.InitialWidth.Value, curErr = expression.Run(nil, w.InitialWidth.Expression, screen).AsFloat()
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial width expression: %w", curErr))
	}
	w.MaxHeight.Value, curErr = expression.Run(nil, w.MaxHeight.Expression, screen).AsFloat()
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving max width expression: %w", curErr))
	}
	w.InitialPositionX.Value, curErr = expression.Run(nil, w.InitialPositionX.Expression, screen).AsFloat()
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial x-position expression: %w", curErr))
	}
	w.InitialZoom.Value, curErr = expression.Run(nil, w.InitialZoom.Expression, screen).AsFloat()
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial zoom expression: %w", curErr))
	}

	w.InitialPositionY.Value, curErr = expression.Run(nil, w.InitialPositionY.Expression, screen).AsFloat()
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial y-position expression: %w", curErr))
	}

	// in case of grow top the axis is inverted (0,0 is left bottom - instead of left top)
	if w.GrowTop {
		w.InitialPositionY.Value = float64(screen.CurrentScreen.Dimension.Height) - w.InitialPositionY.Value
	}

	return
}
