package model

type Help struct {
	Arg    bool `yaml:"help,omitempty" short:"h" usage:"Show this help"`
	Env    bool `yaml:"help-env,omitempty" usage:"Show help for environment variables"`
	Yaml   bool `yaml:"help-config,omitempty" usage:"Show help for config file"`
	Styles bool `yaml:"help-styles,omitempty" usage:"Show help of code styles"`
	Expr   bool `yaml:"help-expression,omitempty" usage:"Show help for expressions"`
	Tool   bool `yaml:"help-tool,omitempty" usage:"Show help for tools"`

	GenYaml  bool `yaml:"config-template,omitempty" usage:"Generate yaml template with default values and all possible options"`
	DumpYaml bool `yaml:"config-dump,omitempty" usage:"Dump the current configuration as yaml"`
}
