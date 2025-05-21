package mistral

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/mistral"
)

type Mistral struct {
	client *mistral.Model
}

func New(opts []mistral.Option) (common.Model, error) {
	result := &Mistral{}

	var err error
	result.client, err = mistral.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating Mistral LLM: %w", err)
	}

	return result, nil
}

func (m *Mistral) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return m.client.Call(ctx, prompt, options...)
}

func (m *Mistral) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return m.client.GenerateContent(ctx, messages, options...)
}

func (m *Mistral) Close() error {
	return nil
}
