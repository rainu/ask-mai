package model

import (
	"github.com/kirsle/configdir"
	"path"
)

type History struct {
	Path string `yaml:"path,omitempty" usage:"The path to the history file. If empty, no history will be used"`
}

func (h *History) SetDefaults() {
	confPath := configdir.LocalConfig("ask-mai")
	confPath = path.Join(confPath, "history")

	h.Path = confPath
}

func (h *History) GetUsage(field string) string {
	switch field {
	}
	return ""
}

func (h *History) Validate() error {
	return nil
}
