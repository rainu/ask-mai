package llm

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIConfig struct {
	APIKey     string `config:"api-key" usage:"API Key"`
	APIType    string `config:"api-type"`
	APIVersion string `config:"api-version" usage:"API Version"`

	Model        string `config:"model" usage:"Model"`
	BaseUrl      string `config:"base-url" usage:"BaseUrl"`
	Organization string `config:"organization" usage:"Organization"`
}

func (c *OpenAIConfig) GetUsage(field string) string {
	switch field {
	case "APIType":
		return fmt.Sprintf("OpenAI API Type (%s, %s, %s)", openai.APITypeOpenAI, openai.APITypeAzure, openai.APITypeAzureAD)
	}
	return ""
}

func configureOpenai(c *OpenAIConfig) {
	flag.StringVar(&c.APIKey, "openai-api-key", "", "OpenAI API Key")
	flag.StringVar(&c.APIType, "openai-api-type", string(openai.APITypeOpenAI), fmt.Sprintf("OpenAI API Type (%s, %s, %s)", openai.APITypeOpenAI, openai.APITypeAzure, openai.APITypeAzureAD))
	flag.StringVar(&c.APIVersion, "openai-api-version", "", "OpenAI API Version")
	flag.StringVar(&c.Model, "openai-model", "gpt-4o-mini", "OpenAI chat model")
	flag.StringVar(&c.BaseUrl, "openai-base-url", "", "OpenAI API Base-URL")
	flag.StringVar(&c.Organization, "openai-organization", "", "OpenAI Organization")
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
