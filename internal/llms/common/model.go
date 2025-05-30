package common

import (
	"github.com/tmc/langchaingo/llms"
)

type Model interface {
	llms.Model

	Close() error
	PatchTools(*[]llms.Tool) error
	ConsumptionOf(*llms.ContentResponse) Consumption
}
