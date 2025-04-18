//go:build debug

package main

import (
	"github.com/rainu/ask-mai/config/model"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	onStartUp = func(c *model.Config) {
		go func() {
			slog.Info("Start pprof server on " + c.Debug.PprofAddress)
			err := http.ListenAndServe(c.Debug.PprofAddress, nil)
			slog.Info("Stop pprof server.", "error", err)
		}()
	}
}
