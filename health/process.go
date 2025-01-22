//go:build !linux

package health

import "context"

func ObserveProcess(ctx context.Context, pid int32, threshold float64, callback func()) {
}
