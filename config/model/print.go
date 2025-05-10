package model

import (
	"fmt"
	"io"
	"os"
)

const (
	PrinterFormatPlain = "plain"
	PrinterFormatJSON  = "json"
	PrinterTargetOut   = "stdout"
	PrinterTargetErr   = "stderr"
)

type PrinterConfig struct {
	Targets    []io.WriteCloser `yaml:"-"`
	TargetsRaw []string         `yaml:"targets,omitempty"`
	Format     string           `yaml:"format,omitempty" short:"f"`
}

func (p *PrinterConfig) SetDefaults() {
	p.Format = PrinterFormatJSON
	p.Targets = []io.WriteCloser{os.Stdout}
	p.TargetsRaw = []string{PrinterTargetOut}
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

func (p *PrinterConfig) Close() {
	for _, target := range p.Targets {
		target.Close()
	}
}
