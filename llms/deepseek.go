package llms

import (
	"fmt"
	"github.com/tmc/langchaingo/llms/openai"
)

func NewDeepSeek(opts []openai.Option) (Model, error) {
	// DeepSeek aims to provide the same rest interface as OpenAI
	// So that is the reason we use the OpenAI implementation here
	result, err := NewOpenAI(opts)
	if err != nil {
		return nil, fmt.Errorf("error creating DeepSeek LLM: %w", err)
	}
	return result, nil
}
