package model

type History struct {
	Path string `yaml:"path" usage:"The path to the history file. If empty, no history will be used"`
}

func (h *History) GetUsage(field string) string {
	switch field {
	}
	return ""
}

func (h *History) Validate() error {
	return nil
}
