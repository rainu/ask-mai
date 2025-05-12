package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/config/model/common"
	"github.com/rainu/ask-mai/expression"
	"github.com/rainu/ask-mai/llms"
	"github.com/rainu/go-yacl"
)

type AnythingLLMConfig struct {
	BaseURL   string                  `yaml:"base-url,omitempty" usage:"Base URL"`
	Token     common.Secret           `yaml:"token,omitempty" usage:"Token"`
	Workspace string                  `yaml:"workspace,omitempty" usage:"Workspace"`
	Thread    AnythingLLMThreadConfig `yaml:"thread,omitempty" usage:"Thread "`
}

type AnythingLLMThreadConfig struct {
	Name   common.StringContainer `yaml:"name,omitempty" usage:"Expression: Name of the newly generated thread"`
	Delete bool                   `yaml:"delete,omitempty" usage:"Delete the thread after the session is closed"`
}

func (a *AnythingLLMThreadConfig) SetDefaults() {
	if a.Name.Expression == nil && a.Name.Value == nil {
		a.Name = common.StringContainer{
			Expression: yacl.P(`'ask-mai - ' + new Date().toISOString()`),
		}
	}
}

func (c *AnythingLLMConfig) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("AnythingLLM Base URL is missing")
	}
	if ce := c.Token.Validate(); ce != nil {
		return fmt.Errorf("AnythingLLM Token is missing: %w", ce)
	}
	if c.Workspace == "" {
		return fmt.Errorf("AnythingLLM Workspace is missing")
	}

	if err := expression.Validate(*c.Thread.Name.Expression); err != nil {
		return fmt.Errorf("Invalid AnythingLLM thread name expression: %w", err)
	}

	return nil
}

func (c *AnythingLLMConfig) BuildLLM() (llms.Model, error) {
	tn, err := expression.Run(nil, *c.Thread.Name.Expression, nil).AsString()
	if err != nil {
		return nil, fmt.Errorf("error calculating thread name: %w", err)
	}

	return llms.NewAnythingLLM(
		c.BaseURL,
		string(c.Token.GetOrPanicWithDefaultTimeout()),
		c.Workspace,
		tn,
		c.Thread.Delete,
	)
}
