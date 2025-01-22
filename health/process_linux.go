//go:build linux

package health

import (
	"context"
	"github.com/shirou/gopsutil/v4/process"
	"log/slog"
	"syscall"
	"time"
)

const observationInterval = 250 * time.Millisecond

func ObserveProcess(ctx context.Context, threshold float64, callback func()) {
	threadId, _, errno := syscall.Syscall(syscall.SYS_GETTID, 0, 0, 0)
	if errno != 0 {
		slog.Error("Error getting thread ID. Process health observation is inactive!", "error", errno)
		return
	}

	p, err := process.NewProcess(int32(threadId))
	if err != nil {
		slog.Error("Unable to get process. Process health observation is inactive!", "pid", threadId, "error", err)
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
				slog.Error("Unable to get process usage.", "pid", threadId, "error", err)
				continue
			}
			if percent > threshold {
				callback()
			}
		}
	}()
}
