package model

import (
	"fmt"
	"os"
)

const (
	ThemeDark   = "dark"
	ThemeLight  = "light"
	ThemeSystem = "system"
)

type UIConfig struct {
	Window       WindowConfig     `yaml:"window,omitempty"`
	Prompt       PromptConfig     `yaml:"prompt,omitempty"`
	FileDialog   FileDialogConfig `yaml:"file-dialog,omitempty"`
	Stream       bool             `yaml:"stream,omitempty" short:"s" usage:"Should the output be streamed"`
	QuitShortcut Shortcut         `yaml:"quit,omitempty" usage:"The shortcut for quitting the application: "`
	Theme        string           `yaml:"theme,omitempty"`
	CodeStyle    string           `yaml:"code-style,omitempty" usage:"The code style to use"`
	Language     string           `yaml:"lang,omitempty" usage:"The language to use"`
}

func (u *UIConfig) SetDefaults() {
	u.QuitShortcut = Shortcut{Binding: []string{"Escape"}}
	u.Theme = ThemeSystem
	u.CodeStyle = "github"
	u.Language = os.Getenv("LANG")
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

	if ve := u.QuitShortcut.Validate(); ve != nil {
		return ve
	}

	if ve := u.Prompt.Validate(); ve != nil {
		return ve
	}

	if ve := u.FileDialog.Validate(); ve != nil {
		return ve
	}

	if ve := u.Window.Validate(); ve != nil {
		return ve
	}

	return nil
}
