package controller

import "github.com/wailsapp/wails/v2/pkg/runtime"

// Runtime functions that are used in the controller package.
// These functions can be overridden in tests or other packages if needed.
var (
	RuntimeEventsEmit = runtime.EventsEmit

	RuntimeScreenGetAll      = runtime.ScreenGetAll
	RuntimeWindowShow        = runtime.WindowShow
	RuntimeWindowGetSize     = runtime.WindowGetSize
	RuntimeWindowSetSize     = runtime.WindowSetSize
	RuntimeWindowSetMaxSize  = runtime.WindowSetMaxSize
	RuntimeWindowSetPosition = runtime.WindowSetPosition
	RuntimeWindowCenter      = runtime.WindowCenter

	RuntimeOpenMultipleFilesDialog = runtime.OpenMultipleFilesDialog

	RuntimeHide = runtime.Hide
	RuntimeQuit = runtime.Quit
)
