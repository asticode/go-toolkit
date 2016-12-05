package time

import (
	"context"
	"errors"
	"time"
)

// Vars
var (
	ErrCancelled = errors.New("cancelled")
)

// Sleep is a cancellable sleep
func Sleep(ctx context.Context, d time.Duration) (err error) {
	for {
		select {
		case <-time.After(d):
			return
		case <-ctx.Done():
			err = ErrCancelled
			return
		}
	}
	return
}
