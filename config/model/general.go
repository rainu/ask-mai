package model

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/rainu/ask-mai/config/model/llm"
)

type Config struct {
	UI UIConfig `yaml:"ui"`

	LLM llm.LLMConfig `config:"" yaml:"llm"`

	Printer PrinterConfig `yaml:"print"`

	History History `yaml:"history"`

	Debug DebugConfig `config:"" yaml:"debug"`

	Config string `config:"config" short:"c" yaml:"-" usage:"Path to the configuration yaml file"`

	ActiveProfile string             `config:"profile" short:"P" usage:"Active profile name"`
	Profiles      map[string]*Config `config:"" yaml:"profiles" usage:"Other configuration profiles. Each profile has the same structure as the main configuration."`
}

func (c *Config) Validate() error {
	if ve := c.Debug.Validate(); ve != nil {
		return ve
	}

	if ve := c.UI.Validate(); ve != nil {
		return ve
	}

	if ve := c.LLM.Validate(); ve != nil {
		return ve
	}

	if ve := c.Printer.Validate(); ve != nil {
		return ve
	}

	if ve := c.History.Validate(); ve != nil {
		return ve
	}

	cWithoutProfiles := *c
	cWithoutProfiles.Profiles = nil
	for profileName, profile := range c.Profiles {
		// merge "default"-config into profile
		err := mergo.Merge(profile, cWithoutProfiles, mergo.WithOverrideEmptySlice)
		if err != nil {
			return fmt.Errorf("Error merging profile '%s': %w", profileName, err)
		}

		if ve := profile.Validate(); ve != nil {
			return fmt.Errorf("Error in profile '%s': %w", profileName, ve)
		}
	}

	return nil
}

func (c *Config) GetActiveProfile() *Config {
	profile, ok := c.Profiles[c.ActiveProfile]
	if !ok {
		return c
	}
	return profile
}
