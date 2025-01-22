//go:build !linux

package health

import "context"

func ObserveProcess(ctx context.Context, threshold float64, callback func()) {
}
