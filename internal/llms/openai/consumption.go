package openai

import (
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/tmc/langchaingo/llms"
)

const (
	tokenKeyCompletion = "CompletionTokens"
	tokenKeyInput      = "PromptTokens"
	tokenKeyReasoning  = "ReasoningTokens"
	tokenKeyCached     = "PromptCachedTokens"
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
		if t, ok := choice.GenerationInfo[tokenKeyCached]; ok {
			result.cached += uint64(t.(int))
			result.input -= uint64(t.(int)) // cached tokens are already part of the input
		}
		if t, ok := choice.GenerationInfo[tokenKeyReasoning]; ok {
			result.reason += uint64(t.(int))
		}
	}

	return result
}

type consumption struct {
	input  uint64
	cached uint64
	output uint64
	reason uint64
}

func (t *consumption) Summary() common.ConsumptionSummary {
	base := common.NewCachedConsumption(t.input, t.cached, t.output)
	base["Reasoning"] = t.reason

	return base
}

func (t *consumption) Add(add common.Consumption) {
	if token, ok := add.(*consumption); ok {
		t.input += token.input
		t.cached += token.cached
		t.output += token.output
		t.reason += token.reason
	}
}
