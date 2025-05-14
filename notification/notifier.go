package notification

import (
	"github.com/gen2brain/beeep"
	"log/slog"
	"os"
	"sync"
)

var iconPath string
var iconWg sync.WaitGroup

func SetNotificationIcon(icon []byte) {
	// check system-wide icon first
	if _, err := os.Stat("/usr/share/pixmaps/ask-mai.png"); err == nil {
		iconPath = "/usr/share/pixmaps/ask-mai.png"
	} else if icon != nil {
		go writeIcon(icon)
	}
}

func writeIcon(icon []byte) {
	iconWg.Add(1)
	defer iconWg.Done()

	f, err := os.CreateTemp("", "*ask-mai.icon.png")
	if err != nil {
		slog.Error("failed to create temp file for notification icon", "err", err)
	}
	defer f.Close()

	if _, err := f.Write(icon); err != nil {
		slog.Error("failed to write notification icon", "err", err)
	}
	iconPath = f.Name()
}

func Notify(title, message string) error {
	iconWg.Wait()

	return beeep.Notify(title, message, iconPath)
}
