package llm

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tmc/langchaingo/llms"
	"gopkg.in/yaml.v3"
	"io"
	"reflect"
	"strings"
)

type ToolsConfig struct {
	RawTools []string `config:"function" yaml:"-" usage:"Function definition (json) to use. See Tool-Usage for more information."`

	Tools map[string]FunctionDefinition `config:"-" yaml:"functions"`
}

type FunctionDefinition struct {
	Name          string `yaml:"-" json:"name" usage:"The name of the function."`
	Description   string `yaml:"description" json:"description" usage:"The description of the function."`
	Parameters    any    `yaml:"parameters" json:"parameters" usage:"The parameter definition of the function."`
	Command       string `yaml:"command" json:"command" usage:"The command to execute."`
	NeedsApproval bool   `yaml:"approval" json:"approval" usage:"Needs user approval to be executed."`
}

func PrintToolsUsage(output io.Writer) {
	fmt.Fprintf(output, "\nTool-Usage:"+
		"\nYou can define many functions that can be used by the LLM."+
		"\nThe functions can be given by argument (JSON), Environment (JSON) or config file (YAML)."+
		"\nThe fields are more or less the same for all three methods:\n")

	table := tablewriter.NewWriter(output)
	table.SetBorder(false)
	table.SetHeader([]string{"Name", "Type", "Usage"})
	table.SetAutoWrapText(false)

	fd := FunctionDefinition{}
	val := reflect.ValueOf(fd)
	typ := reflect.TypeOf(fd)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		t := field.Type.String()
		if strings.HasPrefix(t, "interface") {
			t = "any"
		}
		name := field.Tag.Get("json")
		usage := field.Tag.Get("usage")
		table.Append([]string{name, t, usage})
	}
	table.Render()

	fmt.Fprintf(output, "\nJSON:\n")

	exampleDefs := []FunctionDefinition{
		{
			Name:        "createFile",
			Description: "This function creates a file.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path": map[string]any{
						"type":        "string",
						"description": "The path to the file.",
					},
				},
				"required": []string{"path"},
			},
			Command:       "/usr/bin/touch",
			NeedsApproval: true,
		},
		{
			Name:        "echo",
			Description: "This function echoes a message.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"message": map[string]any{
						"type":        "string",
						"description": "The message to echo.",
					},
				},
				"required": []string{"message"},
			},
			Command:       "/usr/bin/echo",
			NeedsApproval: false,
		},
	}

	fdm := map[string]FunctionDefinition{}
	for _, def := range exampleDefs {
		jsonDef, _ := json.MarshalIndent(def, "", " ")
		fmt.Fprintf(output, "\n%s\n", jsonDef)

		fdm[def.Name] = def
	}

	fmt.Fprintf(output, "\nYAML:\n\n")
	ye := yaml.NewEncoder(output)
	ye.SetIndent(2)
	ye.Encode(map[string]any{
		"llm": map[string]any{
			"tool": ToolsConfig{Tools: fdm},
		},
	})
}

func (t *ToolsConfig) Validate() error {
	for i, tool := range t.RawTools {
		var result FunctionDefinition

		err := json.Unmarshal([]byte(tool), &result)
		if err != nil {
			return fmt.Errorf("Invalid tool definition #%d: %w", i, err)
		}

		t.Tools[result.Name] = result
	}

	for cmd, definition := range t.Tools {
		if definition.Command == "" {
			return fmt.Errorf("Command for tool '%s' is missing", cmd)
		}
	}

	return nil
}

func (t *ToolsConfig) AsOptions() (opts []llms.CallOption) {
	var tools []llms.Tool
	for name, tool := range t.Tools {
		tools = append(tools, llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
		})
	}

	if len(tools) > 0 {
		opts = append(opts, llms.WithTools(tools))
	}
	return
}
