package main

import (
	"embed"
	"github.com/rainu/ask-mai/backend"
	"github.com/rainu/ask-mai/backend/anythingllm"
	"github.com/rainu/ask-mai/backend/copilot"
	"github.com/rainu/ask-mai/backend/openai"
	"github.com/rainu/ask-mai/controller"
	"github.com/rainu/ask-mai/io"
	"log"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v2"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	buildMode := slices.ContainsFunc(os.Environ(), func(s string) bool {
		return strings.HasPrefix(s, "tsprefix=")
	})

	config := ParseConfig()
	if !buildMode {
		if err := config.Validate(); err != nil {
			println(err.Error())
			os.Exit(1)
			return
		}
	}

	slog.SetLogLoggerLevel(slog.Level(config.LogLevel))
	defer func() {
		for _, target := range config.Printer.Targets {
			target.Close()
		}
	}()

	printer := io.MultiResponsePrinter{}
	if config.Printer.Format == PrinterFormatPlain {
		for _, target := range config.Printer.Targets {
			printer.Printers = append(printer.Printers, &io.PlainResponsePrinter{Target: target})
		}
	} else if config.Printer.Format == PrinterFormatJSON {
		for _, target := range config.Printer.Targets {
			printer.Printers = append(printer.Printers, &io.JsonResponsePrinter{Target: target})
		}
	}

	var mainWindowController *controller.Controller
	switch config.Backend {
	case BackendCopilot:
		mainWindowController = controller.New(backend.Builder{
			Build: copilot.NewCopilot,
			Type:  backend.TypeSingleShot,
		}, printer, config.Prompt)
	case BackendOpenAI:
		mainWindowController = controller.New(backend.Builder{
			Build: func() (backend.Handle, error) {
				return openai.NewOpenAI(config.OpenAI.APIKey, config.OpenAI.SystemPrompt)
			},
			Type: backend.TypeMultiShot,
		}, printer, config.Prompt)
	case BackendAnythingLLM:
		mainWindowController = controller.New(backend.Builder{
			Build: func() (backend.Handle, error) {
				return anythingllm.NewAnythingLLM(config.AnythingLLM.BaseURL, config.AnythingLLM.Token, config.AnythingLLM.Workspace)
			},
			Type: backend.TypeMultiShot,
		}, printer, config.Prompt)
	}

	// Create application with options
	err := wails.Run(controller.CreateOptions(controller.CreateOptionsArg{
		Controller: mainWindowController,

		Assets: assets,
		Icon:   icon,
		Width:  int(config.Width),
	}))

	if err != nil {
		log.Fatal(err)
	}
}
