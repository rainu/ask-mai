package localai

import (
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/common"
	iOpenai "github.com/rainu/ask-mai/internal/llms/openai"
	"github.com/tmc/langchaingo/llms/openai"
)

func New(opts []openai.Option) (common.Model, error) {
	// LocalAI aims to provide the same rest interface as OpenAI
	// So that is the reason we use the OpenAI implementation here

	result, err := iOpenai.New(opts)
	if err != nil {
		return nil, fmt.Errorf("error creating LocalAI LLM: %w", err)
	}

	return result, nil
}
