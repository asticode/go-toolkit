package time_test

import (
	"testing"
	stltime "time"

	"github.com/asticode/go-toolkit/time"
	"github.com/stretchr/testify/assert"
)

func TestSleep(t *testing.T) {
	var channelCancel = make(chan bool)
	var err error
	go func() {
		err = time.Sleep(stltime.Minute, channelCancel)
	}()
	channelCancel <- true
	assert.EqualError(t, err, time.ErrCancelled.Error())
}
