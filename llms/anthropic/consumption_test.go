package anthropic

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"testing"
)

func TestAnthropic_ConsumptionOf(t *testing.T) {
	testCases := []struct {
		name           string
		contentResp    *llms.ContentResponse
		expectedTokens consumption
	}{
		{
			name:           "Nil Response",
			expectedTokens: consumption{},
		},
		{
			name: "Empty Response",
			contentResp: &llms.ContentResponse{
				Choices: []*llms.ContentChoice{},
			},
			expectedTokens: consumption{},
		},
		{
			name: "All token types",
			contentResp: &llms.ContentResponse{
				Choices: []*llms.ContentChoice{
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyInput:  42,
							tokenKeyOutput: 123,
						},
					},
				},
			},
			expectedTokens: consumption{
				input:  42,
				output: 123,
			},
		},
		{
			name: "Multiple choices",
			contentResp: &llms.ContentResponse{
				Choices: []*llms.ContentChoice{
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyInput:  10,
							tokenKeyOutput: 20,
						},
					},
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyInput:  15,
							tokenKeyOutput: 25,
						},
					},
				},
			},
			expectedTokens: consumption{
				input:  25,
				output: 45,
			},
		},
		{
			name: "Missing information",
			contentResp: &llms.ContentResponse{
				Choices: []*llms.ContentChoice{
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyOutput: 50,
						},
					},
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyInput: 30,
						},
					},
				},
			},
			expectedTokens: consumption{
				input:  30,
				output: 50,
			},
		},
	}

	openAI := &Anthropic{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := openAI.ConsumptionOf(tc.contentResp)

			tokenResult, ok := result.(*consumption)
			assert.True(t, ok, "invalid return type")
			assert.Equal(t, tc.expectedTokens, *tokenResult)
		})
	}
}
