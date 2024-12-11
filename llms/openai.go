package llms

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAI struct {
	client *openai.LLM
}

func NewOpenAI(opts []openai.Option) (Model, error) {
	result := &OpenAI{}

	var err error
	result.client, err = openai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenAI LLM: %w", err)
	}

	return result, nil
}

func (o *OpenAI) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return o.client.Call(ctx, prompt, options...)
}

func (o *OpenAI) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return o.client.GenerateContent(ctx, messages, options...)
}

func (o *OpenAI) Close() error {
	return nil
}
