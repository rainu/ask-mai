package openai

import (
	"github.com/rainu/ask-mai/llms/common"
	"github.com/tmc/langchaingo/llms"
)

const (
	tokenKeyCompletion = "CompletionTokens"
	tokenKeyInput      = "PromptTokens"
	tokenKeyReasoning  = "ReasoningTokens"
)

func (o *OpenAI) ConsumptionOf(resp *llms.ContentResponse) common.Consumption {
	result := &consumption{}
	if resp == nil {
		return result
	}

	for _, choice := range resp.Choices {
		if t, ok := choice.GenerationInfo[tokenKeyCompletion]; ok {
			result.output += uint64(t.(int))
		}
		if t, ok := choice.GenerationInfo[tokenKeyInput]; ok {
			result.input += uint64(t.(int))
		}
		if t, ok := choice.GenerationInfo[tokenKeyReasoning]; ok {
			result.reason += uint64(t.(int))
		}
	}

	return result
}

type consumption struct {
	input  uint64
	output uint64
	reason uint64
}

func (t *consumption) Summary() common.ConsumptionSummary {
	base := common.NewSimpleConsumption(t.input, t.output)
	base["Reasoning"] = t.reason

	return base
}

func (t *consumption) Add(add common.Consumption) {
	if token, ok := add.(*consumption); ok {
		t.input += token.input
		t.output += token.output
		t.reason += token.reason
	}
}
