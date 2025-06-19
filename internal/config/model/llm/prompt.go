package llm

import (
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

type PromptConfig struct {
	System string `yaml:"system,omitempty" short:"S" usage:"System Prompt"`

	InitMessages    []Message `yaml:"init-message,omitempty" usage:"initial message(s) to use: "`
	InitValue       string    `yaml:"init-value,omitempty" short:"p" usage:"the initial prompt to use"`
	InitAttachments []string  `yaml:"init-attachment,omitempty" short:"a" usage:"the initial attachment(s) to use"`
}

func (p *PromptConfig) Validate() error {
	for _, message := range p.InitMessages {
		if ve := message.Validate(); ve != nil {
			return ve
		}
	}
	return nil
}

func (p *PromptConfig) HasInitialValues() bool {
	return p.InitValue != "" || len(p.InitMessages) > 0 || len(p.InitAttachments) > 0
}

type Message struct {
	Role    llms.ChatMessageType `yaml:"role,omitempty"`
	Content string               `yaml:"content,omitempty" usage:"content"`
}

func (m *Message) GetUsage(field string) string {
	switch field {
	case "Role":
		return fmt.Sprintf("role (%s, %s)", llms.ChatMessageTypeHuman, llms.ChatMessageTypeAI)
	}
	return ""
}

func (m *Message) Validate() error {
	if m.Role != llms.ChatMessageTypeHuman && m.Role != llms.ChatMessageTypeAI {
		return fmt.Errorf("Invalid message role '%s', must be one of %s, %s!", m.Role, llms.ChatMessageTypeHuman, llms.ChatMessageTypeAI)
	}
	if m.Content == "" {
		return fmt.Errorf("Message content is required!")
	}
	return nil
}
