package main

import (
	"github.com/rainu/ask-mai/backend"
	"github.com/rainu/ask-mai/backend/anythingllm"
	"github.com/rainu/ask-mai/backend/copilot"
	"github.com/rainu/ask-mai/backend/openai"
	"github.com/rainu/ask-mai/io"
	errView "github.com/rainu/ask-mai/ui/view/error"
	questionView "github.com/rainu/ask-mai/ui/view/question"
	"log/slog"
	"os"
)

func main() {
	config := ParseConfig()
	if err := config.Validate(); err != nil {
		errorWindow := errView.NewWindow(err.Error())
		errorWindow.ShowAndRun()
		os.Exit(1)
		return
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

	var mainWindowController *questionView.Controller
	switch config.Backend {
	case BackendCopilot:
		mainWindowController = questionView.NewController(backend.Builder{
			Build: copilot.NewCopilot,
			Type:  backend.TypeSingleShot,
		}, printer)
	case BackendOpenAI:
		mainWindowController = questionView.NewController(backend.Builder{
			Build: func() (backend.Handle, error) {
				return openai.NewOpenAI(config.OpenAI.APIKey, config.OpenAI.SystemPrompt)
			},
			Type: backend.TypeMultiShot,
		}, printer)
	case BackendAnythingLLM:
		mainWindowController = questionView.NewController(backend.Builder{
			Build: func() (backend.Handle, error) {
				return anythingllm.NewAnythingLLM(config.AnythingLLM.BaseURL, config.AnythingLLM.Token, config.AnythingLLM.Workspace)
			},
			Type: backend.TypeMultiShot,
		}, printer)
	}
	defer mainWindowController.Close()

	mainWindow := questionView.NewWindow(mainWindowController, config.Prompt, config.Width, config.Height)
	mainWindow.ShowAndRun()
}
