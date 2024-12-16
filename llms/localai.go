package llms

import (
	"fmt"
	"github.com/tmc/langchaingo/llms/openai"
)

func NewLocalAI(opts []openai.Option) (Model, error) {
	// LocalAI aims to provide the same rest interface as OpenAI
	// So that is the reason we use the OpenAI implementation here

	result, err := NewOpenAI(opts)
	if err != nil {
		return nil, fmt.Errorf("error creating LocalAI LLM: %w", err)
	}

	return result, nil
}
