package anthropic

import (
	"github.com/rainu/ask-mai/llms/common"
	"github.com/tmc/langchaingo/llms"
)

const (
	tokenKeyInput  = "InputTokens"
	tokenKeyOutput = "OutputTokens"
)

func (a *Anthropic) ConsumptionOf(resp *llms.ContentResponse) common.Consumption {
	result := &consumption{}
	if resp == nil {
		return result
	}

	for _, choice := range resp.Choices {
		if t, ok := choice.GenerationInfo[tokenKeyInput]; ok {
			result.input += uint64(t.(int))
		}
		if t, ok := choice.GenerationInfo[tokenKeyOutput]; ok {
			result.output += uint64(t.(int))
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
