package config

import (
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model"
	"github.com/rainu/go-yacl"
	"os"
	"strings"
)

const EnvironmentPrefix = "ASK_MAI_"

func Parse(arguments []string, env []string) *model.Config {
	// for possible config file path
	cf := &model.ConfigFile{}
	config := yacl.NewConfig(cf, yacl.WithPrefixEnv(EnvironmentPrefix))
	handleErr(func() error { return config.ParseEnvironment(env...) })
	handleErr(func() error { return config.ParseArguments(arguments...) })

	c := &model.Config{}
	config = yacl.NewConfig(c, yacl.WithPrefixEnv(EnvironmentPrefix))
	config.ApplyDefaults()

	processYamlFiles(config, cf.Path)
	handleErr(func() error { return config.ParseEnvironment(env...) })
	handleErr(func() error { return config.ParseArguments(arguments...) })
	checkHelp(c, config)

	if &c.MainProfile != c.GetActiveProfile() {
		// apply arguments to the active profile too
		config = yacl.NewConfig(c.GetActiveProfile(), yacl.WithPrefixEnv(EnvironmentPrefix))
		handleErr(func() error { return config.ParseEnvironment(env...) })
		handleErr(func() error { return config.ParseArguments(arguments...) })
	}

	c.MainProfile.Printer.Targets = nil
	for _, target := range c.MainProfile.Printer.TargetsRaw {
		target = strings.TrimSpace(target)

		if target == model.PrinterTargetOut {
			c.MainProfile.Printer.Targets = append(c.MainProfile.Printer.Targets, os.Stdout)
		} else if target == model.PrinterTargetErr {
			c.MainProfile.Printer.Targets = append(c.MainProfile.Printer.Targets, os.Stderr)
		} else {
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
			if err != nil {
				panic(fmt.Errorf("Error creating printer target file: %w", err))
			}
			c.MainProfile.Printer.Targets = append(c.MainProfile.Printer.Targets, file)
		}
	}

	return c
}

func handleErr(f func() error) {
	if err := f(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
