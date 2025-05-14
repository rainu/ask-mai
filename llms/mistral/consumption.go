package mistral

import (
	"github.com/gage-technologies/mistral-go"
	"github.com/rainu/ask-mai/llms/common"
	"github.com/tmc/langchaingo/llms"
)

const (
	generalInfoKeyUsage = "usage"
)

func (m *Mistral) ConsumptionOf(resp *llms.ContentResponse) common.Consumption {
	result := &consumption{}
	if resp == nil {
		return result
	}

	for _, choice := range resp.Choices {
		if m, ok := choice.GenerationInfo[generalInfoKeyUsage]; ok {
			if usage, ok := m.(mistral.UsageInfo); ok {
				result.input += uint64(usage.PromptTokens)
				result.output += uint64(usage.CompletionTokens)
			}
		}
	}

	return result
}

type consumption struct {
	input  uint64
	output uint64
}

func (t *consumption) Summary() common.ConsumptionSummary {
	return common.NewSimpleConsumption(t.input, t.output)
}

func (t *consumption) Add(add common.Consumption) {
	if token, ok := add.(*consumption); ok {
		t.input += token.input
		t.output += token.output
	}
}
