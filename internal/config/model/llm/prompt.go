package llm

import (
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools"
	"github.com/rainu/ask-mai/internal/mcp/client"
	"github.com/tmc/langchaingo/llms"
	"log/slog"
)

type PromptConfig struct {
	System string `yaml:"system,omitempty" short:"S" usage:"System Prompt"`

	InitMessages    []Message  `yaml:"init-message,omitempty" usage:"initial message(s) to use: "`
	InitToolCalls   []ToolCall `yaml:"init-tool-call,omitempty" usage:"initial tool call(s) to use: "`
	InitValue       string     `yaml:"init-value,omitempty" short:"p" usage:"the initial prompt to use"`
	InitAttachments []string   `yaml:"init-attachment,omitempty" short:"a" usage:"the initial attachment(s) to use"`
}

func (p *PromptConfig) Validate() error {
	for _, message := range p.InitMessages {
		if ve := message.Validate(); ve != nil {
			return ve
		}
	}
	for _, toolCall := range p.InitToolCalls {
		if ve := toolCall.Validate(); ve != nil {
			return ve
		}
	}
	return nil
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

const (
	ToolCallTypeBuiltin = "_builtin"
	ToolCallTypeCustom  = "_custom"
)

type ToolCall struct {
	Server    string         `yaml:"server,omitempty"`
	Name      string         `yaml:"name,omitempty" usage:"name of the tool"`
	Arguments map[string]any `yaml:"args,omitempty" usage:"arguments to pass"`
}

func (t *ToolCall) GetUsage(field string) string {
	switch field {
	case "Server":
		return fmt.Sprintf("server which proiveds the tool. '%s' for builtin tools, '%s' for any custom function or any configured mcp-server.", ToolCallTypeBuiltin, ToolCallTypeCustom)
	}
	return ""
}

func (t *ToolCall) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("Tool name is required!")
	}
	if t.GetArguments() == "" {
		return fmt.Errorf("Tool arguments can not be marshalled to JSON, please check the arguments!")
	}

	return nil
}

func (t *ToolCall) GetArguments() string {
	if len(t.Arguments) == 0 {
		return "{}"
	}

	result, _ := json.Marshal(t.Arguments)
	return string(result)
}

func (t *ToolCall) GetTransporter(c *tools.Config) client.Transporter {
	switch t.Server {
	case ToolCallTypeBuiltin:
		return c.GetBuiltinTransporter()
	case ToolCallTypeCustom:
		return c.GetCustomTransporter()
	}

	mcpServer, ok := c.McpServer[t.Server]
	if !ok {
		slog.Warn("Invalid initial tool call config! Unknown MCP-Server!", slog.String("server", t.Server))
		return nil
	}
	return &mcpServer
}

func (t *ToolCall) GetUniqName() string {
	serverName := t.Server
	if serverName == ToolCallTypeBuiltin {
		serverName = tools.ServerNameBuiltin
	} else if serverName == ToolCallTypeCustom {
		serverName = tools.ServerNameCustom
	}

	return tools.GetUniqToolName(serverName, t.Name)
}
