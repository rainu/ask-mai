package llm

import (
	"fmt"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIConfig struct {
	APIKey     string `yaml:"api-key" usage:"API Key"`
	APIType    string `yaml:"api-type"`
	APIVersion string `yaml:"api-version" usage:"API Version"`

	Model        string `yaml:"model" usage:"Model"`
	BaseUrl      string `yaml:"base-url" usage:"BaseUrl"`
	Organization string `yaml:"organization" usage:"Organization"`
}

func (c *OpenAIConfig) GetUsage(field string) string {
	switch field {
	case "APIType":
		return fmt.Sprintf("OpenAI API Type (%s, %s, %s)", openai.APITypeOpenAI, openai.APITypeAzure, openai.APITypeAzureAD)
	}
	return ""
}

func (c *OpenAIConfig) AsOptions() (opts []openai.Option) {
	if c.APIKey != "" {
		opts = append(opts, openai.WithToken(c.APIKey))
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
	if c.APIKey == "" {
		return fmt.Errorf("OpenAI API Key is missing")
	}
	if c.APIType != "" && c.APIType != string(openai.APITypeOpenAI) && c.APIType != string(openai.APITypeAzure) && c.APIType != string(openai.APITypeAzureAD) {
		return fmt.Errorf("OpenAI API Type is invalid")
	}

	return nil
}
