package sync_test

import (
	"testing"
	"time"

	"github.com/asticode/go-toolkit/sync"
	"github.com/rs/xlog"
	"github.com/stretchr/testify/assert"
)

func TestRWMutex_IsDeadlocked(t *testing.T) {
	var m = sync.NewRWMutex(xlog.NopLogger, "test")
	d, _ := m.IsDeadlocked(time.Millisecond)
	assert.False(t, d)
	m.Lock()
	d, c := m.IsDeadlocked(time.Millisecond)
	assert.True(t, d)
	assert.Contains(t, c, "github.com/asticode/go-toolkit/sync/mutex_test.go:16")
}
