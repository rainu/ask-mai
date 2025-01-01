package config

import "fmt"

type FileDialogConfig struct {
	DefaultDirectory           string `config:"default-dir" usage:"The default directory for the file dialog"`
	ShowHiddenFiles            bool   `config:"show-hidden" usage:"Should the file dialog show hidden files"`
	CanCreateDirectories       bool   `config:"can-create-dirs" usage:"Should the file dialog be able to create directories"`
	ResolvesAliases            bool   `config:"resolves-aliases" usage:"Should the file dialog resolve aliases"`
	TreatPackagesAsDirectories bool   `config:"treat-packages-as-dirs" usage:"Should the file dialog treat packages as directories"`

	FilterDisplay []string `config:"filter-display" usage:"The filter display names for the file dialog. For example: \"Image Files (*.jpg, *.png)\""`
	FilterPattern []string `config:"filter-pattern" usage:"The filter patterns for the file dialog. For example: \"*.jpg;*.png\""`
}

func (c *FileDialogConfig) Validate() error {
	if len(c.FilterDisplay) != len(c.FilterPattern) {
		return fmt.Errorf("Invalid file dialog filter configuration: it must have the same number of display names and patterns")
	}

	return nil
}
