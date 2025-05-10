package config

import (
	"fmt"
	"github.com/rainu/go-yacl"
	"io"
	"log/slog"
	"os"
	"path"
)

var yamlLookupLocations = func() (result []string) {
	result = append(result, "/"+path.Join("etc", ".ask-mai.yml"))
	result = append(result, "/"+path.Join("etc", ".ask-mai.yaml"))
	result = append(result, "/"+path.Join("etc", "ask-mai", "config.yml"))
	result = append(result, "/"+path.Join("etc", "ask-mai", "config.yaml"))
	result = append(result, "/"+path.Join("usr", "local", "etc", ".ask-mai.yml"))
	result = append(result, "/"+path.Join("usr", "local", "etc", ".ask-mai.yaml"))
	result = append(result, "/"+path.Join("usr", "local", "etc", "ask-mai", "config.yml"))
	result = append(result, "/"+path.Join("usr", "local", "etc", "ask-mai", "config.yaml"))

	if home, err := os.UserHomeDir(); err == nil {
		result = append(result, path.Join(home, ".ask-mai.yml"))
		result = append(result, path.Join(home, ".ask-mai.yaml"))
		result = append(result, path.Join(home, ".config", ".ask-mai.yml"))
		result = append(result, path.Join(home, ".config", ".ask-mai.yaml"))
		result = append(result, path.Join(home, ".config", "ask-mai", "config.yml"))
		result = append(result, path.Join(home, ".config", "ask-mai", "config.yaml"))
	}

	binDir := path.Dir(os.Args[0])
	result = append(result, path.Join(binDir, ".ask-mai.yml"))
	result = append(result, path.Join(binDir, ".ask-mai.yaml"))

	if wd, err := os.Getwd(); err == nil {
		result = append(result, path.Join(wd, ".ask-mai.yml"))
		result = append(result, path.Join(wd, ".ask-mai.yaml"))
	}

	return
}

func processYamlFiles(config *yacl.Config, configFilePath string) {
	for _, location := range yamlLookupLocations() {
		processYamlFile(config, location)
	}
	if configFilePath != "" {
		processYamlFile(config, configFilePath)
	}
}

func processYamlFile(config *yacl.Config, path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	slog.Debug("Processing yaml file", "file", path)
	err = processYaml(config, f)
	if err != nil {
		panic(fmt.Errorf("unable to process yaml file %s: %w", path, err))
	}
}

func processYaml(config *yacl.Config, source io.Reader) error {
	err := config.ParseYaml(source)
	if err != nil {
		return fmt.Errorf("error while decoding yaml: %w", err)
	}
	return nil
}
