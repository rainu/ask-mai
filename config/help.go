package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/rainu/ask-mai/config/llm/tools"
	"github.com/rainu/ask-mai/expression"
	"github.com/rainu/ask-mai/llms/tools/command"
	http2 "github.com/rainu/ask-mai/llms/tools/http"
	flag "github.com/spf13/pflag"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

type NoopLooger struct {
}

func (n NoopLooger) Print(message string) {
}

func (n NoopLooger) Trace(message string) {
}

func (n NoopLooger) Debug(message string) {
}

func (n NoopLooger) Info(message string) {
}

func (n NoopLooger) Warning(message string) {
}

func (n NoopLooger) Error(message string) {
}

func (n NoopLooger) Fatal(message string) {
}

func printHelpArgs(output io.Writer, fields resolvedFieldInfos) {
	fmt.Fprintf(output, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func printHelpEnv(output io.Writer, fields resolvedFieldInfos) {
	fmt.Fprintf(output, "Available environment variables:\n")

	table := tablewriter.NewWriter(output)
	table.SetBorder(false)
	table.SetHeader([]string{"Name", "Usage"})
	table.SetAutoWrapText(false)

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Env < fields[j].Env
	})
	for _, field := range fields {
		env := field.Env
		if strings.HasPrefix(field.Value.Type().String(), "*[]") {
			env += "_N"
		}
		table.Append([]string{env, field.Usage})
	}
	table.Render()
}

func printHelpConfig(output io.Writer, fields resolvedFieldInfos) {
	sort.Slice(fields, func(i, j int) bool {
		return strings.Join(fields[i].YamlPath, "") < strings.Join(fields[j].YamlPath, "")
	})

	fmt.Fprintf(output, "Yaml keys:\n")

	table := tablewriter.NewWriter(output)
	table.SetBorder(false)
	table.SetHeader([]string{"Name", "Usage"})
	table.SetAutoWrapText(false)

	for _, field := range fields {
		yamlKey := strings.TrimLeft(strings.Join(field.YamlPath, "."), ".")
		if strings.HasSuffix(yamlKey, "-") {
			continue
		}
		table.Append([]string{yamlKey, field.Usage})
	}
	table.Render()

	fmt.Fprintf(output, "\nYaml lookup file locations:\n")
	for _, location := range yamlLookupLocations() {
		fmt.Fprintf(output, "  - %s\n", location)
	}
}

func printHelpStyles(output io.Writer) {
	fmt.Fprintf(output, "\nAvailable code styles:\n")
	for _, style := range availableCodeStyles {
		fmt.Fprintf(output, "  - %s\n", style)
	}
}

func printHelpExpression(output io.Writer) {
	fmt.Fprintf(output, "The expression language is JavaScript. You can use the following variables and functions:\n")
	fmt.Fprintf(output, "\nFunctions:\n")
	fmt.Fprintf(output, "  - %s: writes a message to the console.\n", expression.FuncNameLog)

	fmt.Fprintf(output, "\nVariables:\n")

	wails.Run(&options.App{
		StartHidden: true,
		Frameless:   true,
		Width:       1,
		Height:      1,
		OnStartup: func(ctx context.Context) {
			screens, err := runtime.ScreenGetAll(ctx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}

			fmt.Fprintf(output, "  const %s = ", expression.VarNameScreens)

			variables := expression.SetScreens(screens)
			je := json.NewEncoder(output)
			je.SetIndent("  ", "  ")
			je.Encode(variables)
		},
		OnDomReady: func(ctx context.Context) {
			runtime.Quit(ctx)
		},
		AssetServer: &assetserver.Options{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		},
		Logger:           &NoopLooger{},
		WindowStartState: options.Minimised,
	})
}

func printHelpTool(output io.Writer) {
	fmt.Fprintf(output, "Tool-Usage:"+
		"\nYou can define many functions that can be used by the LLM."+
		"\nThe functions can be given by argument (JSON), Environment (JSON) or config file (YAML)."+
		"\nThe fields are more or less the same for all three methods:\n")

	table := tablewriter.NewWriter(output)
	table.SetBorder(false)
	table.SetHeader([]string{"Name", "Type", "Usage"})
	table.SetAutoWrapText(false)

	fields := scanConfigTags(nil, &tools.FunctionDefinition{})
	for _, field := range fields {
		t := strings.TrimPrefix(field.Value.Type().String(), "*")
		if strings.HasPrefix(t, "interface") {
			t = "any"
		}
		table.Append([]string{strings.Replace(field.Flag, ",omitempty", "", -1), t, field.Usage})
	}
	table.Render()

	exampleDefs := []tools.FunctionDefinition{
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
			Command: "/usr/bin/touch",
			Environment: map[string]string{
				"USER":  "rainu",
				"SHELL": "/bin/bash",
			},
			WorkingDir:    "/tmp",
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
			AdditionalEnvironment: map[string]string{
				"ASK_MAI_ARGS": "$@",
			},
			Command:       "/usr/bin/echo",
			NeedsApproval: false,
		},
	}

	fmt.Fprintf(output, "\nJSON:\n")

	fdm := map[string]tools.FunctionDefinition{}
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
			"tool": map[string]any{
				"functions": fdm,
			},
		},
	})

	fmt.Fprintf(output, "\nThe LLM will respond the arguments as JSON. You can use the following placeholders in the command:\n")
	fmt.Fprintf(output, "  - $@: all arguments (1:1 the JSON from the LLM)\n")
	fmt.Fprintf(output, "  - $<varName>: the value of <varName> in the LLM's JSON\n")
	fmt.Fprintf(output, "\nExamples:\n")

	table = tablewriter.NewWriter(output)
	table.SetBorder(false)
	table.SetHeader([]string{"Pattern", "LLM's JSON", "Result"})
	table.SetAutoWrapText(false)

	table.Append([]string{`/usr/bin/echo $@`, `{"message": "hello world"}`, `/usr/bin/echo {"message": "hello world"}`})
	table.Append([]string{`/usr/bin/echo $message`, `{"message": "hello world"}`, `/usr/bin/echo hello world`})
	table.Append([]string{`/usr/bin/echo "$message"`, `{"message": "hello world"}`, `/usr/bin/echo "hello world"`})
	table.Append([]string{`/usr/bin/echo "$message"`, `{}`, `/usr/bin/echo ""`})

	table.Render()

	fmt.Fprintf(output, "\nYou can also use these placeholder in (additional) environment and working directory variables.\n")

	fmt.Fprintf(output, "\nIt is also possible to define a JavaScript expression (file).\n")
	fmt.Fprintf(output, "You can use the following variables and functions:\n")
	fmt.Fprintf(output, "\nFunctions:\n")
	fmt.Fprintf(output, "  - %s(...args): writes a message to the console.\n", expression.FuncNameLog)

	js := bytes.Buffer{}
	je := json.NewEncoder(&js)
	je.SetIndent("   ", "  ")
	je.Encode(command.CommandDescriptor{
		Command:   "/path/to/command",
		Arguments: []string{"arg1", "...argN"},
		Environment: map[string]string{
			"ENV_VAR": "value",
		},
		AdditionalEnvironment: map[string]string{
			"ADDITIONAL_ENV_VAR": "value",
		},
		WorkingDirectory: "/path/to/working/dir",
	})
	fmt.Fprintf(output, "  - %s(%s): run a command.\n", expression.FuncNameRun, strings.TrimSpace(js.String()))

	js = bytes.Buffer{}
	je = json.NewEncoder(&js)
	je.SetIndent("   ", "  ")
	je.Encode(http2.CallDescriptor{
		Method: http.MethodPost,
		Url:    "https://example.com",
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		StringBody: `{"msg": "hello world"}`,
	})
	fmt.Fprintf(output, "  - %s(%s): do a http call.\n", expression.FuncNameFetch, strings.TrimSpace(js.String()))

	fmt.Fprintf(output, "\nVariables:\n")
	fmt.Fprintf(output, "  const %s = ", expression.VarNameContext)

	je = json.NewEncoder(output)
	je.SetIndent("  ", "  ")
	je.Encode(tools.CommandVariables{
		FunctionDefinition: tools.FunctionDefinition{
			Name:        "jsEcho",
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
			CommandExpr:   fmt.Sprintf(`"Echo: " + JSON.parse(%s.args).message`, expression.VarNameContext),
			NeedsApproval: false,
		},
		Arguments: `{"message": "hello world"}`,
	})

	fmt.Fprintf(output, "\nJavaScript command expression examples:")
	fmt.Fprintf(output, "\n\n  // parse llm's JSON, run the command and return its result\n")
	fmt.Fprintf(output, `
  const pa = JSON.parse(v.args)
  const cmdDescriptor = {
   "command": "echo",
   "arguments": ["Echo:", pa.message]
  }
  
  `+expression.FuncNameRun+`(cmdDescriptor)`)

	fmt.Fprintf(output, "\n\n  // catches possible execution error\n")
	fmt.Fprintf(output, `
  let result = ""
  try {
     result = `+expression.FuncNameRun+`({ "command": "__doesNotExists__" })
  } catch (e) {
     result = "Error: " + e
  }
  result
`)

}
