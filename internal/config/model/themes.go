package model

import "fmt"

type Themes struct {
	Dark   *Theme           `yaml:"dark,omitempty" json:"dark,omitempty" usage:"Dark theme "`
	Light  *Theme           `yaml:"light,omitempty" json:"light,omitempty" usage:"Light theme "`
	Custom map[string]Theme `yaml:"custom,omitempty" json:"custom,omitempty" usage:"Custom themes (See https://vuetifyjs.com/en/features/theme/) "`
}

type Theme struct {
	Dark      bool              `yaml:"-" json:"dark"`
	Colors    map[string]string `yaml:"colors" json:"colors" usage:"colors"`
	Variables map[string]string `yaml:"variables" json:"variables" usage:"variables"`
}

func (t *Themes) SetDefaults() {
	if t.Dark == nil {
		t.Dark = &Theme{
			Dark: true,
			Colors: map[string]string{
				"chat-system-message":    "",
				"chat-assistant-message": "#E0E0E0",
				"chat-tool-call":         "#E0E0E0",
				"chat-user-message":      "#69F0AE",
			},
		}
		t.Light = &Theme{
			Dark: false,
			Colors: map[string]string{
				"chat-system-message":    "",
				"chat-assistant-message": "#E0E0E0",
				"chat-tool-call":         "#E0E0E0",
				"chat-user-message":      "#69F0AE",
			},
		}
	}
}

func (t *Themes) Validate() error {
	for _, theme := range []string{ThemeDark, ThemeLight, ThemeSystem} {
		if _, exists := t.Custom[theme]; exists {
			return fmt.Errorf("The theme '%s' is reserved for the system theme and cannot be defined in the configuration.", theme)
		}
	}

	dt := Themes{}
	dt.SetDefaults()

	setDefaultColors(t.Dark, dt.Dark.Colors)
	setDefaultColors(t.Light, dt.Light.Colors)
	for tn := range t.Custom {
		tv := t.Custom[tn]
		setDefaultColors(&tv, dt.Dark.Colors)
		t.Custom[tn] = tv
	}
	t.Dark.Dark = true
	t.Light.Dark = false

	return nil
}

func setDefaultColors(t *Theme, colors map[string]string) {
	for dColorName, dColorValue := range colors {
		if _, exists := t.Colors[dColorName]; !exists {
			t.Colors[dColorName] = dColorValue
		}
	}
}
