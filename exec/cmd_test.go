package exec_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/asticode/go-toolkit/exec"
	"github.com/stretchr/testify/assert"
)

func TestWithTimeout(t *testing.T) {
	// Init
	cmd := exec.NewCmd("sleep", "0.5")

	// Success
	cmd.Timeout = 1 * time.Second
	assert.Equal(t, "sleep 0.5", cmd.String())
	_, _, err := exec.Exec(cmd)
	assert.NoError(t, err)

	// Timeout
	cmd.Timeout = time.Millisecond
	_, _, err = exec.Exec(cmd)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("Timeout of %v reached", cmd.Timeout), err.Error()[:22])

	// Cancel
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _, err = exec.Exec(cmd)
	}()
	cmd.ChannelCancel <- true
	wg.Wait()
	assert.Error(t, err)
	assert.Equal(t, "Command was cancelled, no process to kill", err.Error()[:41])
}
