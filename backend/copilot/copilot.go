package copilot

import (
	"bufio"
	"context"
	"fmt"
	"github.com/rainu/ask-mai/backend"
	cmdchain "github.com/rainu/go-command-chain"
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

type Copilot struct {
	ctx      context.Context
	cancelFn context.CancelFunc
}

type interaction struct {
	Prefix string
	Output string
}

func NewCopilot() (backend.Handle, error) {
	return &Copilot{}, nil
}

func (c *Copilot) AskSomething(chat []backend.Message) (string, error) {
	return c.AskSomethingWithContext(context.Background(), chat)
}

func (c *Copilot) AskSomethingWithContext(ctx context.Context, chat []backend.Message) (string, error) {
	if len(chat) == 0 {
		return "", fmt.Errorf("empty history provided")
	}

	c.ctx, c.cancelFn = context.WithCancel(ctx)

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
			JoinWithContext(c.ctx, "gh", "copilot", "suggest", chat[len(chat)-1].Content, "--target", "shell", "--shell-out", tempFile.Name()).
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
	if c.cancelFn != nil {
		c.cancelFn()
	}

	return nil
}
