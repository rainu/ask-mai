package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/config/expression"
	flag "github.com/spf13/pflag"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func printUsage(output io.Writer, fields resolvedFieldInfos) {
	fmt.Fprintf(output, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

	fmt.Fprintf(output, "\nAvailable environment variables:\n")

	maxLen := 0
	sort.Slice(fields, func(i, j int) bool {
		//the sort algorithm must iterate all elements
		//so here we "recycle" the loop to determine max-length
		if len(fields[i].Env) > maxLen {
			maxLen = len(fields[i].Env)
		}
		return fields[i].Env < fields[j].Env
	})
	for _, field := range fields {
		env := field.Env
		if strings.HasPrefix(field.Value.Type().String(), "*[]") {
			env += "_N"
		}
		fmt.Fprintf(output, "  %s%s\t%s\n", env, strings.Repeat(" ", maxLen-len(env)), field.Usage)
	}

	sort.Slice(fields, func(i, j int) bool {
		return strings.Join(fields[i].YamlPath, "") < strings.Join(fields[j].YamlPath, "")
	})

	fmt.Fprintf(output, "\nYaml keys:\n")
	for _, field := range fields {
		yamlKey := strings.TrimLeft(strings.Join(field.YamlPath, "."), ".")
		if strings.HasSuffix(yamlKey, "-") {
			continue
		}
		fmt.Fprintf(output, "  %s%s\t%s\n", yamlKey, strings.Repeat(" ", maxLen-len(yamlKey)), field.Usage)
	}

	fmt.Fprintf(output, "\nYaml lookup file locations:\n")
	for _, location := range yamlLookupLocations() {
		fmt.Fprintf(output, "  - %s\n", location)
	}

	fmt.Fprintf(output, "\nAvailable code styles:\n")
	for _, style := range availableCodeStyles {
		fmt.Fprintf(output, "  - %s\n", style)
	}

	fmt.Fprintf(output, "\nThe expression language is JavaScript. You can use the following variables and functions:\n")
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

			fmt.Fprintf(output, "  const %s = ", expression.VarNameVariables)

			variables := expression.FromScreens(screens)
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
