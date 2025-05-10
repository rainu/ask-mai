package model

import (
	"fmt"
	"strings"
)

type Shortcut struct {
	Binding []string `yaml:"binding,omitempty" usage:"The binding for the shortcut"`
	Code    []string `yaml:"-"`
	Alt     []bool   `yaml:"-"`
	Ctrl    []bool   `yaml:"-"`
	Meta    []bool   `yaml:"-"`
	Shift   []bool   `yaml:"-"`
}

func (s *Shortcut) Validate() error {
	s.Code = make([]string, len(s.Binding))
	s.Alt = make([]bool, len(s.Binding))
	s.Ctrl = make([]bool, len(s.Binding))
	s.Meta = make([]bool, len(s.Binding))
	s.Shift = make([]bool, len(s.Binding))

	for b := range s.Binding {
		normalized := strings.Replace(strings.ToLower(s.Binding[b]), " ", "", -1)
		parts := strings.Split(normalized, "+")
		for _, part := range parts {
			switch part {
			case "alt":
				s.Alt[b] = true
			case "ctrl":
				s.Ctrl[b] = true
			case "meta":
				s.Meta[b] = true
			case "shift":
				s.Shift[b] = true
			default:
				s.Code[b] = part
			}
		}
		if s.Code[b] == "" {
			return fmt.Errorf("Invalid binding '%s': no key code", s.Binding[b])
		}
	}

	return nil
}
