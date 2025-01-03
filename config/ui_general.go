package config

import (
	"fmt"
)

const (
	ThemeDark   = "dark"
	ThemeLight  = "light"
	ThemeSystem = "system"
)

type UIConfig struct {
	Window       WindowConfig     `yaml:"window"`
	Prompt       PromptConfig     `yaml:"prompt"`
	FileDialog   FileDialogConfig `yaml:"file-dialog"`
	Stream       bool             `yaml:"stream" short:"s" usage:"Should the output be streamed"`
	QuitShortcut Shortcut         `yaml:"quit" usage:"The shortcut for quitting the application: "`
	Theme        string           `yaml:"theme"`
	CodeStyle    string           `yaml:"code-style" usage:"The code style to use"`
	Language     string           `yaml:"lang" usage:"The language to use"`
}

func (u *UIConfig) GetUsage(field string) string {
	switch field {
	case "Theme":
		return fmt.Sprintf("The theme to use ('%s', '%s', '%s')", ThemeLight, ThemeDark, ThemeSystem)
	}
	return ""
}

func (u *UIConfig) Validate() error {
	if u.Theme != ThemeDark && u.Theme != ThemeLight && u.Theme != ThemeSystem {
		return fmt.Errorf("Invalid theme")
	}

	if ve := u.FileDialog.Validate(); ve != nil {
		return ve
	}

	if ve := u.Window.Validate(); ve != nil {
		return ve
	}

	return nil
}
