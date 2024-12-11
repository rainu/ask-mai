package llms

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type Anthropic struct {
	client *anthropic.LLM
}

func NewAnthropic(opts []anthropic.Option) (Model, error) {
	result := &Anthropic{}

	var err error
	result.client, err = anthropic.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating Anthropic LLM: %w", err)
	}

	return result, nil
}

func (o *Anthropic) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return o.client.Call(ctx, prompt, options...)
}

func (o *Anthropic) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return o.client.GenerateContent(ctx, messages, options...)
}

func (o *Anthropic) Close() error {
	return nil
}
