package config

type Shortcut struct {
	Code  string `config:"key" usage:"key-code"`
	Alt   bool   `config:"alt" usage:"alt-key must be pressed"`
	Ctrl  bool   `config:"ctrl" usage:"control-key must be pressed"`
	Meta  bool   `config:"meta" usage:"meta-key must be pressed"`
	Shift bool   `config:"shift" usage:"shift-key must be pressed"`
}
