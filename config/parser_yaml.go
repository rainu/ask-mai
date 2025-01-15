package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log/slog"
	"os"
	"path"
)

func yamlLookupLocations() (result []string) {
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

func processYamlFiles(c *Config) {
	for _, location := range yamlLookupLocations() {
		processYamlFile(location, c)
	}
	if c.Config != "" {
		processYamlFile(c.Config, c)
	}
}

func processYamlFile(path string, c *Config) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	slog.Debug("Processing yaml file", "file", path)
	err = processYaml(f, c)
	if err != nil {
		panic(fmt.Errorf("unable to process yaml file %s: %w", path, err))
	}
}

func processYaml(source io.Reader, c *Config) error {
	err := yaml.NewDecoder(source).Decode(c)
	if err != nil {
		return fmt.Errorf("error while decoding yaml: %w", err)
	}
	return nil
}
