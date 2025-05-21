package model

import (
	"fmt"
	"github.com/rainu/go-yacl"
)

type FileDialogConfig struct {
	DefaultDirectory           string `yaml:"default-dir,omitempty" usage:"The default directory for the file dialog"`
	ShowHiddenFiles            *bool  `yaml:"show-hidden,omitempty" usage:"Should the file dialog show hidden files"`
	CanCreateDirectories       *bool  `yaml:"can-create-dirs,omitempty" usage:"Should the file dialog be able to create directories"`
	ResolveAliases             *bool  `yaml:"resolve-aliases,omitempty" usage:"Should the file dialog resolve aliases"`
	TreatPackagesAsDirectories *bool  `yaml:"treat-packages-as-dirs,omitempty" usage:"Should the file dialog treat packages as directories"`

	FilterDisplay []string `yaml:"filter-display,omitempty" usage:"The filter display names for the file dialog. For example: \"Image Files (*.jpg, *.png)\""`
	FilterPattern []string `yaml:"filter-pattern,omitempty" usage:"The filter patterns for the file dialog. For example: \"*.jpg;*.png\""`
}

func (c *FileDialogConfig) SetDefaults() {
	if c.ShowHiddenFiles == nil {
		c.ShowHiddenFiles = yacl.P(true)
	}
	if c.CanCreateDirectories == nil {
		c.CanCreateDirectories = yacl.P(false)
	}
	if c.ResolveAliases == nil {
		c.ResolveAliases = yacl.P(true)
	}
	if c.TreatPackagesAsDirectories == nil {
		c.TreatPackagesAsDirectories = yacl.P(true)
	}
}

func (c *FileDialogConfig) Validate() error {
	if len(c.FilterDisplay) != len(c.FilterPattern) {
		return fmt.Errorf("Invalid file dialog filter configuration: it must have the same number of display names and patterns")
	}

	return nil
}
