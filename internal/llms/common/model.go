package common

import (
	"github.com/tmc/langchaingo/llms"
)

type Model interface {
	llms.Model

	Close() error
	ConsumptionOf(*llms.ContentResponse) Consumption
}
