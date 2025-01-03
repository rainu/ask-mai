package config

import (
	"fmt"
	"io"
)

const (
	PrinterFormatPlain = "plain"
	PrinterFormatJSON  = "json"
	PrinterTargetOut   = "stdout"
	PrinterTargetErr   = "stderr"
)

type PrinterConfig struct {
	Format     string           `yaml:"format" short:"f"`
	Targets    []io.WriteCloser `yaml:"-"`
	TargetsRaw []string         `yaml:"targets"`
}

func (p *PrinterConfig) GetUsage(field string) string {
	switch field {
	case "Format":
		return fmt.Sprintf("Response printer format (%s, %s)", PrinterFormatPlain, PrinterFormatJSON)
	case "TargetsRaw":
		return fmt.Sprintf("Response printer targets (%s, %s, <path/to/file>)", PrinterTargetOut, PrinterTargetErr)
	}
	return ""
}

func (p *PrinterConfig) Validate() error {
	if p.Format != PrinterFormatJSON && p.Format != PrinterFormatPlain {
		return fmt.Errorf("Invalid response printer format")
	}

	return nil
}
