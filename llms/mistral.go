package llms

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/mistral"
)

type Mistral struct {
	client *mistral.Model
}

func NewMistral(opts []mistral.Option) (Model, error) {
	result := &Mistral{}

	var err error
	result.client, err = mistral.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating Mistral LLM: %w", err)
	}

	return result, nil
}

func (o *Mistral) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return o.client.Call(ctx, prompt, options...)
}

func (o *Mistral) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return o.client.GenerateContent(ctx, messages, options...)
}

func (o *Mistral) Close() error {
	return nil
}
