package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model/common"
	llmCommon "github.com/rainu/ask-mai/internal/llms/common"
	"github.com/rainu/ask-mai/internal/llms/google"
	"github.com/rainu/go-yacl"
	"github.com/tmc/langchaingo/llms/googleai"
	"time"
)

type GoogleAIConfig struct {
	APIKey        common.Secret `yaml:"api-key,omitempty" usage:"API Key"`
	Model         string        `yaml:"model,omitempty" usage:"Model"`
	HarmThreshold *int32        `yaml:"harm-threshold,omitempty"`

	ToolCacheTTL *time.Duration `yaml:"tool-cache-ttl,omitempty" usage:"TTL for tool cache. 0 means no caching. Minimum is 1 minute."`
}

func (c *GoogleAIConfig) SetDefaults() {
	if c.Model == "" {
		c.Model = "gemini-2.0-flash-lite"
	}
	if c.HarmThreshold == nil {
		c.HarmThreshold = yacl.P(int32(googleai.HarmBlockUnspecified))
	}
	if c.ToolCacheTTL == nil {
		c.ToolCacheTTL = yacl.P(5 * time.Minute)
	}
}

func (c *GoogleAIConfig) GetUsage(field string) string {
	switch field {
	case "HarmThreshold":
		return fmt.Sprintf("The safety/harm setting for the model, potentially limiting any harmful content it may generate"+
			"\n\t\t%d - threshold is unspecified"+
			"\n\t\t%d - content with NEGLIGIBLE will be allowed"+
			"\n\t\t%d - content with NEGLIGIBLE and LOW will be allowed"+
			"\n\t\t%d - content with NEGLIGIBLE, LOW, and MEDIUM will be allowed"+
			"\n\t\t%d - all content will be allowed", googleai.HarmBlockUnspecified, googleai.HarmBlockLowAndAbove, googleai.HarmBlockMediumAndAbove, googleai.HarmBlockOnlyHigh, googleai.HarmBlockNone)
	}
	return ""
}

func (c *GoogleAIConfig) AsOptions() (opts []googleai.Option) {
	if v := c.APIKey.GetOrPanicWithDefaultTimeout(); v != nil {
		opts = append(opts, googleai.WithAPIKey(string(v)))
	}
	if c.Model != "" {
		opts = append(opts, googleai.WithDefaultModel(c.Model))
	}
	if c.HarmThreshold != nil {
		opts = append(opts, googleai.WithHarmThreshold(googleai.HarmBlockThreshold(*c.HarmThreshold)))
	}

	return
}

func (c *GoogleAIConfig) Validate() error {
	if ce := c.APIKey.Validate(); ce != nil {
		return fmt.Errorf("GoogleAI API Key is missing: %w", ce)
	}
	if c.HarmThreshold != nil {
		switch googleai.HarmBlockThreshold(*c.HarmThreshold) {
		case googleai.HarmBlockUnspecified:
		case googleai.HarmBlockLowAndAbove:
		case googleai.HarmBlockMediumAndAbove:
		case googleai.HarmBlockOnlyHigh:
		case googleai.HarmBlockNone: // valid values
		default:
			return fmt.Errorf("Invalid harm threshold value: %d", c.HarmThreshold)
		}
	}
	if c.ToolCacheTTL != nil && *c.ToolCacheTTL > 0 && *c.ToolCacheTTL < time.Minute {
		return fmt.Errorf("Tool cache TTL must be at least 1 minute, got: %s", c.ToolCacheTTL)
	}

	return nil
}

func (c *GoogleAIConfig) BuildLLM() (llmCommon.Model, error) {
	return google.New(c.AsOptions(), *c.ToolCacheTTL)
}
