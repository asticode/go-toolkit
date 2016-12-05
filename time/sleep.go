package time

import (
	"errors"
	"time"
)

// Vars
var (
	ErrCancelled = errors.New("cancelled")
)

// Sleep is a cancellable sleep
func Sleep(d time.Duration, channelCancel chan bool) (err error) {
	for {
		select {
		case <-time.After(d):
			return
		case <-channelCancel:
			err = ErrCancelled
			return
		}
	}
	return
}
