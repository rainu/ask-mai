//go:build linux

package health

import (
	"context"
	"github.com/shirou/gopsutil/v4/process"
	"log/slog"
	"time"
)

const observationInterval = 250 * time.Millisecond

func ObserveProcess(ctx context.Context, pid int32, threshold float64, callback func()) {
	p, err := process.NewProcess(pid)
	if err != nil {
		slog.Error("Unable to get process. Observation is inactive!", "pid", pid, "error", err)
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(observationInterval):
			}

			percent, err := p.Percent(observationInterval)
			if err != nil {
				slog.Error("Unable to get process usage.", "pid", pid, "error", err)
				continue
			}
			if percent > threshold {
				callback()
			}
		}
	}()
}
