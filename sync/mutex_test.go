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
	d := m.IsDeadlocked(time.Millisecond)
	assert.False(t, d)
	m.Lock()
	d = m.IsDeadlocked(time.Millisecond)
	assert.True(t, d)
	var s = m.LastSuccessfulLockCaller()
	assert.Equal(t, 16, s.Line)
	assert.Equal(t, "github.com/asticode/go-toolkit/sync_test.TestRWMutex_IsDeadlocked", s.Function)
}
