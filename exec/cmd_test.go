package exec_test

import (
	"sync"
	"testing"
	"time"

	"github.com/asticode/go-toolkit/exec"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestWithTimeout(t *testing.T) {
	// Success
	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	var cmd = exec.NewCmd(ctx, "sleep", "0.5")
	assert.Equal(t, "sleep 0.5", cmd.String())
	_, _, err := exec.Exec(cmd)
	assert.NoError(t, err)

	// Timeout
	ctx, _ = context.WithTimeout(context.Background(), time.Millisecond)
	cmd = exec.NewCmd(ctx, "sleep", "0.5")
	_, _, err = exec.Exec(cmd)
	assert.EqualError(t, err, "signal: killed")

	// Cancel
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	cmd = exec.NewCmd(ctx, "sleep", "0.5")
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _, err = exec.Exec(cmd)
	}()
	cancel()
	wg.Wait()
	assert.EqualError(t, err, "context canceled")
}
