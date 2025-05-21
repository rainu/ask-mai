package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

func (c *Controller) Log(level string, jsonArgs string) {
	var args []any
	json.Unmarshal([]byte(jsonArgs), &args)

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("[console.%s]", level))
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			msg.WriteString(fmt.Sprintf(" %v", args[i]))
		}
		args = nil
	} else {
		args = []any{"args", jsonArgs}
	}

	switch level {
	case "info":
		slog.Info(msg.String(), args...)
	case "warn":
		slog.Warn(msg.String(), args...)
	case "error":
		slog.Error(msg.String(), args...)
	default:
		slog.Debug(msg.String(), args...)
	}
}
