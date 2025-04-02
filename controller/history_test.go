package controller

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"testing"
)

func TestController_HistoryConverting(t *testing.T) {
	toTest := LLMMessages{
		{
			Role: string(llms.ChatMessageTypeHuman),
			ContentParts: []LLMMessageContentPart{{
				Type:    LLMMessageContentPartTypeText,
				Content: "How much time is it?",
			}},
		},
		{
			Role: string(llms.ChatMessageTypeTool),
			ContentParts: []LLMMessageContentPart{{
				Type: LLMMessageContentPartTypeToolCall,
				Call: LLMMessageCall{
					Id:        "123",
					Function:  "__getTime",
					Arguments: "{}",
					Result: &LLMMessageCallResult{
						Content:    "2025-04-02 11:09:14.506479929 +0200 CEST m=+6.615035201",
						Error:      "",
						DurationMs: 13,
					},
				},
			}},
		},
		{
			Role: string(llms.ChatMessageTypeAI),
			ContentParts: []LLMMessageContentPart{{
				Type:    LLMMessageContentPartTypeText,
				Content: "It is 11:09 a.m. on April 2, 2025.",
			}},
		},
	}

	entry := historyMessagesToEntry(toTest)
	result := historyEntry2Messages(entry)

	assert.Equal(t, toTest, result)
}
