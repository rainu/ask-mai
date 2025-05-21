package openai

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"testing"
)

func TestOpenAI_ConsumptionOf(t *testing.T) {
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
							tokenKeyInput:      42,
							tokenKeyCompletion: 123,
							tokenKeyReasoning:  7,
						},
					},
				},
			},
			expectedTokens: consumption{
				input:  42,
				output: 123,
				reason: 7,
			},
		},
		{
			name: "Multiple choices",
			contentResp: &llms.ContentResponse{
				Choices: []*llms.ContentChoice{
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyInput:      10,
							tokenKeyCompletion: 20,
							tokenKeyReasoning:  5,
						},
					},
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyInput:      15,
							tokenKeyCompletion: 25,
							tokenKeyReasoning:  10,
						},
					},
				},
			},
			expectedTokens: consumption{
				input:  25,
				output: 45,
				reason: 15,
			},
		},
		{
			name: "Missing information",
			contentResp: &llms.ContentResponse{
				Choices: []*llms.ContentChoice{
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyCompletion: 50,
						},
					},
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyInput: 30,
						},
					},
					{
						GenerationInfo: map[string]interface{}{
							tokenKeyReasoning: 20,
						},
					},
				},
			},
			expectedTokens: consumption{
				input:  30,
				output: 50,
				reason: 20,
			},
		},
	}

	openAI := &OpenAI{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := openAI.ConsumptionOf(tc.contentResp)

			tokenResult, ok := result.(*consumption)
			assert.True(t, ok, "invalid return type")
			assert.Equal(t, tc.expectedTokens, *tokenResult)
		})
	}
}
