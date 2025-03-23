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
	resp, err := o.client.GenerateContent(ctx, messages, options...)
	if err != nil {
		return resp, err
	}

	if len(resp.Choices) > 1 {
		// handle tool call -
		// instead of put tool calls into the same choice as text generating (such like openai does),
		// the anthropic implementation will put tool calls in separate choices!
		// To streamline the process, here we will put the tool calls in the first message of the response.

		for i := len(resp.Choices) - 1; i > 0; i-- {
			if len(resp.Choices[i].ToolCalls) > 0 {
				resp.Choices[0].ToolCalls = append(resp.Choices[0].ToolCalls, resp.Choices[i].ToolCalls...)
			}

			// remove the tool call choice
			resp.Choices = resp.Choices[:i]
		}
	}

	return resp, err
}

func (o *Anthropic) Close() error {
	return nil
}
