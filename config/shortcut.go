package config

type Shortcut struct {
	Code  string `yaml:"key" usage:"key-code"`
	Alt   bool   `yaml:"alt" usage:"alt-key must be pressed"`
	Ctrl  bool   `yaml:"ctrl" usage:"control-key must be pressed"`
	Meta  bool   `yaml:"meta" usage:"meta-key must be pressed"`
	Shift bool   `yaml:"shift" usage:"shift-key must be pressed"`
}
