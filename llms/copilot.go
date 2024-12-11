package llms

import (
	"bufio"
	"context"
	"fmt"
	cmdchain "github.com/rainu/go-command-chain"
	"github.com/tmc/langchaingo/llms"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"
)

var interruptionPattern = []string{
	"? What would you like the shell command to do?",
}

type Copilot struct{}

type interaction struct {
	Prefix string
	Output string
}

func NewCopilot() (Model, error) {
	return &Copilot{}, nil
}

func (c *Copilot) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("empty messages provided")
	}

	prompt := ""
	for _, part := range messages[len(messages)-1].Parts {
		if textPart, ok := part.(llms.TextContent); ok {
			prompt += textPart.Text
		}
	}

	result, err := c.Call(ctx, prompt, options...)
	if err != nil {
		return nil, fmt.Errorf("error calling copilot: %w", err)
	}

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{Content: result},
		},
	}, nil
}

func (c *Copilot) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	inputIn, inputOut := io.Pipe()
	outputIn, outputOut := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(2)

	tempFile, err := os.CreateTemp("", "Copilot-*.txt")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	var commandErr error
	go func() {
		defer wg.Done()
		defer outputOut.Close()

		commandErr = cmdchain.Builder().WithInput(inputIn).
			JoinWithContext(ctx, "gh", "copilot", "suggest", prompt, "--target", "shell", "--shell-out", tempFile.Name()).
			Finalize().WithOutput(outputOut).
			Run()
	}()

	interactions := []interaction{
		{Prefix: "? Select an option", Output: "execute\n"},
		{Prefix: "? Are you sure you want to execute the suggested command?", Output: "yes\n"},
	}

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(outputIn)
		for scanner.Scan() {
			line := scanner.Text()
			slog.Debug(line)

			finishConversation := false

			if strings.HasPrefix(line, interactions[0].Prefix) {
				_, err := io.WriteString(inputOut, interactions[0].Output)
				if err != nil {
					slog.Error("Error writing to stdin:", err)
					continue
				}

				if len(interactions) > 1 {
					interactions = interactions[1:]
				} else {
					finishConversation = true
				}
			} else {
				finishConversation = slices.ContainsFunc(interruptionPattern, func(pattern string) bool {
					return strings.HasPrefix(line, pattern)
				})
			}

			if finishConversation {
				inputOut.Close()
			}
		}
	}()

	wg.Wait()

	if commandErr != nil {
		return "", fmt.Errorf("error running command: %w", commandErr)
	}

	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("error reading temporary file: %w", err)
	}
	return string(content), nil
}

func (c *Copilot) Close() error {
	return nil
}
