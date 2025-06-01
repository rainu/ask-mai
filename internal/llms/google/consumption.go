package google

import (
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/tmc/langchaingo/llms"
)

const (
	tokenKeyInput  = "input_tokens"
	tokenKeyCached = "cached_tokens"
	tokenKeyOutput = "output_tokens"
)

func (g *Google) ConsumptionOf(resp *llms.ContentResponse) common.Consumption {
	result := &consumption{}
	if resp == nil {
		return result
	}

	for _, choice := range resp.Choices {
		if t, ok := choice.GenerationInfo[tokenKeyOutput]; ok {
			result.output += uint64(t.(int32))
		}
		if t, ok := choice.GenerationInfo[tokenKeyCached]; ok {
			result.cached += uint64(t.(int32))

			// input tokens included the cached token count too
			// so we need to subtract them here
			result.input -= uint64(t.(int32))
		}
		if t, ok := choice.GenerationInfo[tokenKeyInput]; ok {
			result.input += uint64(t.(int32))
		}
	}

	return result
}

type consumption struct {
	input  uint64
	cached uint64
	output uint64
}

func (t *consumption) Summary() common.ConsumptionSummary {
	base := common.NewCachedConsumption(t.input, t.cached, t.output)

	return base
}

func (t *consumption) Add(add common.Consumption) {
	if token, ok := add.(*consumption); ok {
		t.input += token.input
		t.cached += token.cached
		t.output += token.output
	}
}
