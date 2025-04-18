package model

type Profile struct {
	Active string `config:"profile" yaml:"active" short:"P" usage:"The active profile name"`

	Icon        string `config:"" yaml:"icon" usage:"Icon to use for the profile"`
	Description string `config:"" yaml:"description" usage:"Description of the profile"`
}
