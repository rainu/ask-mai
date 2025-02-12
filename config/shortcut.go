package config

import (
	"fmt"
	"strings"
)

type Shortcut struct {
	Binding string `yaml:"binding" usage:"The binding for the shortcut"`
	Code    string `config:"-"`
	Alt     bool   `config:"-"`
	Ctrl    bool   `config:"-"`
	Meta    bool   `config:"-"`
	Shift   bool   `config:"-"`
}

func (s *Shortcut) Validate() error {
	normalized := strings.Replace(strings.ToLower(s.Binding), " ", "", -1)
	parts := strings.Split(normalized, "+")
	for _, part := range parts {
		switch part {
		case "alt":
			s.Alt = true
		case "ctrl":
			s.Ctrl = true
		case "meta":
			s.Meta = true
		case "shift":
			s.Shift = true
		default:
			s.Code = part
		}
	}
	if s.Code == "" {
		return fmt.Errorf("Invalid binding '%s': no key code", s.Binding)
	}

	return nil
}
