package model

import (
	"fmt"
)

type Config struct {
	ConfigFile `yaml:",inline"`

	MainProfile Profile     `yaml:",inline,omitempty"`
	DebugConfig DebugConfig `yaml:",inline,omitempty"`

	ActiveProfile string              `yaml:"active-profile,omitempty" short:"P" usage:"The active profile name"`
	Profiles      map[string]*Profile `yaml:"profiles,omitempty" usage:"Configuration profiles. Each profile has the same structure as the main configuration: "`
	Themes        Themes              `yaml:"themes,omitempty" usage:"Theme settings for the application: "`

	Version bool `yaml:"version,omitempty" short:"v" usage:"Show the version"`

	Help Help `yaml:",inline,omitempty"`
}

type ConfigFile struct {
	Path string `yaml:"config-file,omitempty" short:"c" usage:"Path to the configuration yaml file"`
}

func (c *Config) Validate() error {
	if ve := c.MainProfile.Validate(); ve != nil {
		return ve
	}
	if ve := c.DebugConfig.Validate(); ve != nil {
		return ve
	}
	if ve := c.Themes.Validate(); ve != nil {
		return ve
	}

	for profileName, profile := range c.Profiles {
		if ve := profile.Validate(); ve != nil {
			return fmt.Errorf("Error in profile '%s': %w", profileName, ve)
		}
		if profile.UI.Theme != "" {
			if profile.UI.Theme != ThemeDark && profile.UI.Theme != ThemeLight && profile.UI.Theme != ThemeSystem {
				if _, exists := c.Themes.Custom[profile.UI.Theme]; !exists {
					return fmt.Errorf("Unknown theme: %s", profile.UI.Theme)
				}
			}
		}
	}

	return nil
}

func (c *Config) GetActiveProfile() *Profile {
	profile, ok := c.Profiles[c.ActiveProfile]
	if !ok {
		return &c.MainProfile
	}
	return profile
}

func (c *Config) GetProfiles() map[string]ProfileMeta {
	result := map[string]ProfileMeta{}
	result[""] = c.MainProfile.Meta

	for name, config := range c.Profiles {
		result[name] = config.Meta
	}

	return result
}
