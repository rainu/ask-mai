package config

import (
	"fmt"
)

const (
	ThemeDark   = "dark"
	ThemeLight  = "light"
	ThemeSystem = "system"

	MinMaxPositionLeft  = "left"
	MinMaxPositionRight = "right"
	MinMaxPositionNone  = "none"
)

type UIConfig struct {
	Window         WindowConfig     `yaml:"window"`
	Prompt         PromptConfig     `yaml:"prompt"`
	FileDialog     FileDialogConfig `yaml:"file-dialog"`
	Stream         bool             `yaml:"stream" short:"s" usage:"Should the output be streamed"`
	QuitShortcut   Shortcut         `yaml:"quit" usage:"The shortcut for quitting the application: "`
	Theme          string           `yaml:"theme"`
	MinMaxPosition string           `yaml:"min-max-button-position"`
	CodeStyle      string           `yaml:"code-style" usage:"The code style to use"`
	Language       string           `yaml:"lang" usage:"The language to use"`
}

type MinMaxConfig struct {
	Position string `yaml:"position"`
}

func (u *UIConfig) GetUsage(field string) string {
	switch field {
	case "Theme":
		return fmt.Sprintf("The theme to use ('%s', '%s', '%s')", ThemeLight, ThemeDark, ThemeSystem)
	case "MinMaxPosition":
		return fmt.Sprintf("The position of the min/max button ('%s', '%s', '%s')", MinMaxPositionLeft, MinMaxPositionRight, MinMaxPositionNone)
	}
	return ""
}

func (u *UIConfig) Validate() error {
	if u.Theme != ThemeDark && u.Theme != ThemeLight && u.Theme != ThemeSystem {
		return fmt.Errorf("Invalid theme")
	}
	if u.MinMaxPosition != MinMaxPositionLeft && u.MinMaxPosition != MinMaxPositionRight && u.MinMaxPosition != MinMaxPositionNone {
		return fmt.Errorf("Invalid min-max position")
	}

	if ve := u.FileDialog.Validate(); ve != nil {
		return ve
	}

	if ve := u.Window.Validate(); ve != nil {
		return ve
	}

	return nil
}
