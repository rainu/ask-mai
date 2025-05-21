package copilot

import (
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/tmc/langchaingo/llms"
)

func (c *Copilot) ConsumptionOf(resp *llms.ContentResponse) common.Consumption {
	return &common.UnknownConsumption{}
}
