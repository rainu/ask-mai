package anthropic

import (
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/tmc/langchaingo/llms"
)

const (
	tokenKeyInput          = "InputTokens"
	tokenKeyOutput         = "OutputTokens"
	tokenKeyCachedCreation = "CacheCreationInputTokens"
	tokenKeyCachedRead     = "CacheReadInputTokens"
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
		if t, ok := choice.GenerationInfo[tokenKeyCachedCreation]; ok {
			result.cachedCreation += uint64(t.(int))
		}
		if t, ok := choice.GenerationInfo[tokenKeyCachedRead]; ok {
			result.cachedRead += uint64(t.(int))
		}
	}

	return result
}

type consumption struct {
	input  uint64
	output uint64

	cachedCreation uint64
	cachedRead     uint64
}

func (t *consumption) Summary() common.ConsumptionSummary {
	base := common.NewCachedConsumption(t.input, t.cachedRead, t.output)
	base["CacheCreation"] = t.cachedCreation

	return base
}

func (t *consumption) Add(add common.Consumption) {
	if token, ok := add.(*consumption); ok {
		t.input += token.input
		t.output += token.output
		t.cachedCreation += token.cachedCreation
		t.cachedRead += token.cachedRead
	}
}
