package main

import (
	"github.com/rainu/ask-mai/backend"
	"github.com/rainu/ask-mai/backend/anythingllm"
	"github.com/rainu/ask-mai/backend/copilot"
	"github.com/rainu/ask-mai/backend/openai"
	errView "github.com/rainu/ask-mai/ui/view/error"
	questionView "github.com/rainu/ask-mai/ui/view/question"
	"log/slog"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	config := ParseConfig()
	if err := config.Validate(); err != nil {
		errorWindow := errView.NewWindow(err.Error())
		errorWindow.ShowAndRun()
		return
	}

	var mainWindowController *questionView.Controller
	switch config.Backend {
	case BackendCopilot:
		mainWindowController = questionView.NewController(backend.Builder{
			Build: copilot.NewCopilot,
			Type:  backend.TypeSingleShot,
		})
	case BackendOpenAI:
		mainWindowController = questionView.NewController(backend.Builder{
			Build: func() (backend.Handle, error) {
				return openai.NewOpenAI(config.OpenAI.APIKey)
			},
			Type: backend.TypeMultiShot,
		})
	case BackendAnythingLLM:
		mainWindowController = questionView.NewController(backend.Builder{
			Build: func() (backend.Handle, error) {
				return anythingllm.NewAnythingLLM(config.AnythingLLM.BaseURL, config.AnythingLLM.Token, config.AnythingLLM.Workspace)
			},
			Type: backend.TypeMultiShot,
		})
	}
	defer mainWindowController.Close()

	mainWindow := questionView.NewWindow(mainWindowController, config.Prompt)
	mainWindow.ShowAndRun()
}
