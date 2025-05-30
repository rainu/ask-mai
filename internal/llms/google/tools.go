package google

import (
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tmc/langchaingo/llms"
)

func (g *Google) PatchTools(tools *[]llms.Tool) error {
	// the GoogleAI's implementation needs a map as tool-parameter
	for i := range *tools {
		inputSchema, ok := (*tools)[i].Function.Parameters.(mcp.ToolInputSchema)
		if !ok {
			continue
		}

		var err error
		(*tools)[i].Function.Parameters, err = convertToolParameters(inputSchema)
		if err != nil {
			return fmt.Errorf("error converting tool parameters for %s: %w", (*tools)[i].Function.Name, err)
		}
	}

	return nil
}

func convertToolParameters(inputSchema mcp.ToolInputSchema) (map[string]any, error) {
	result := map[string]any{}
	raw, err := json.Marshal(inputSchema)
	if err != nil {
		return nil, fmt.Errorf("error marshalling tool input schema: %w", err)
	}
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling tool input schema: %w", err)
	}

	return result, nil
}
