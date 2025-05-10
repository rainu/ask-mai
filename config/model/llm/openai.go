package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/config/model/common"
	"github.com/rainu/ask-mai/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIConfig struct {
	APIKey     common.Secret `yaml:"api-key,omitempty" usage:"API Key"`
	APIType    string        `yaml:"api-type,omitempty"`
	APIVersion string        `yaml:"api-version,omitempty" usage:"API Version"`

	Model        string `yaml:"model,omitempty" usage:"Model"`
	BaseUrl      string `yaml:"base-url,omitempty" usage:"BaseUrl"`
	Organization string `yaml:"organization,omitempty" usage:"Organization"`
}

func (c *OpenAIConfig) SetDefaults() {
	c.APIType = string(openai.APITypeOpenAI)
	c.Model = "gpt-4o-mini"
}

func (c *OpenAIConfig) GetUsage(field string) string {
	switch field {
	case "APIType":
		return fmt.Sprintf("OpenAI API Type (%s, %s, %s)", openai.APITypeOpenAI, openai.APITypeAzure, openai.APITypeAzureAD)
	}
	return ""
}

func (c *OpenAIConfig) AsOptions() (opts []openai.Option) {
	if v := c.APIKey.GetOrPanicWithDefaultTimeout(); v != nil {
		opts = append(opts, openai.WithToken(string(v)))
	}
	if c.APIType != "" {
		opts = append(opts, openai.WithAPIType(openai.APIType(c.APIType)))
	}
	if c.APIVersion != "" {
		opts = append(opts, openai.WithAPIVersion(c.APIVersion))
	}
	if c.Model != "" {
		opts = append(opts, openai.WithModel(c.Model))
	}
	if c.BaseUrl != "" {
		opts = append(opts, openai.WithBaseURL(c.BaseUrl))
	}
	if c.Organization != "" {
		opts = append(opts, openai.WithOrganization(c.Organization))
	}

	return
}

func (c *OpenAIConfig) Validate() error {
	if ce := c.APIKey.Validate(); ce != nil {
		return fmt.Errorf("OpenAI API Key is missing: %w", ce)
	}
	if c.APIType != "" && c.APIType != string(openai.APITypeOpenAI) && c.APIType != string(openai.APITypeAzure) && c.APIType != string(openai.APITypeAzureAD) {
		return fmt.Errorf("OpenAI API Type is invalid")
	}

	return nil
}

func (c *OpenAIConfig) BuildLLM() (llms.Model, error) {
	return llms.NewOpenAI(c.AsOptions())
}
