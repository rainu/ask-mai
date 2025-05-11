package model

import (
	"dario.cat/mergo"
	"fmt"
)

type Config struct {
	ConfigFile `yaml:",inline"`

	MainProfile Profile     `yaml:",inline,omitempty"`
	DebugConfig DebugConfig `yaml:",inline,omitempty"`

	ActiveProfile string              `yaml:"active-profile,omitempty" short:"P" usage:"The active profile name"`
	Profiles      map[string]*Profile `yaml:"profiles,omitempty" usage:"Configuration profiles. Each profile has the same structure as the main configuration: "`

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

	for profileName, profile := range c.Profiles {
		// merge mainProfile into current profile
		err := mergo.Merge(profile, &c.MainProfile, mergo.WithOverrideEmptySlice, mergo.WithoutDereference)
		if err != nil {
			return fmt.Errorf("Error merging profile '%s': %w", profileName, err)
		}

		if ve := profile.Validate(); ve != nil {
			return fmt.Errorf("Error in profile '%s': %w", profileName, ve)
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
