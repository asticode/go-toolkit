package toolkit_test

import (
	"sync"
	"testing"

	"github.com/asticode/go-toolkit"
	"github.com/stretchr/testify/assert"
)

func TestCanceller_Cancel(t *testing.T) {
	var c = toolkit.NewCanceller()
	var ch1, ch2 = c.NewChannel(), c.NewChannel()
	defer c.Close(ch1)
	defer c.Close(ch2)
	var wg = &sync.WaitGroup{}
	wg.Add(2)
	var count int
	go func() {
		for {
			select {
			case <-ch1:
				count += 1
				wg.Done()
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-ch2:
				count += 2
				wg.Done()
				return
			}
		}
	}()
	c.Cancel()
	wg.Wait()
	assert.Equal(t, 3, count)
}

func TestCanceller_Reset(t *testing.T) {
	var c = toolkit.NewCanceller()
	c.Cancel()
	c.Reset()
	assert.False(t, c.Cancelled())
}
