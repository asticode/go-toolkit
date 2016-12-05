package sync_test

import (
	"testing"
	"time"

	"github.com/asticode/go-toolkit/sync"
	"github.com/rs/xlog"
	"github.com/stretchr/testify/assert"
)

func mockStack() []byte {
	return []byte(`goroutine 1 [running]:
	runtime/debug.Stack(0xc42012e280, 0xc4200f8e17, 0xc4200f8d78)
	/usr/local/go/src/runtime/debug/stack.go:24 +0x79
	github.com/asticode/go-sync.glob..func2(0x0, 0x0, 0xc4200f8da8)
	/home/asticode/projects/go/src/github.com/asticode/go-sync/sync.go:55 +0x22
	github.com/asticode/go-sync.(*RWMutex).LogMessage(0xc420105e00, 0x7a2108, 0x11, 0x7d0710, 0xc4200f8e60)
	/home/asticode/projects/go/src/github.com/asticode/go-sync/sync.go:70 +0x3b
	github.com/asticode/go-sync.(*RWMutex).RUnlock.func1(0xc42001d980, 0xc4201dc000, 0x14, 0xc420105e00)
	/home/asticode/projects/go/src/github.com/asticode/go-sync/sync.go:147 +0x70
	github.com/asticode/go-sync.(*RWMutex).RUnlock(0xc420105e00, 0xc42001d980, 0xc4201dc000, 0x14)
	/home/asticode/projects/go/src/github.com/asticode/go-sync/sync.go:151 +0x3d
	main.(*Worker).Retire(0xc4200f1800, 0x90cfa0, 0xc420018d70)
	/home/asticode/projects/go/src/github.com/asticode/myproject/worker.go:174 +0x11d
	main.main()
	/home/asticode/projects/go/src/github.com/asticode/myproject/main.go:76 +0x571`)
}

func TestRWMutex_Parent(t *testing.T) {
	sync.Stack = func() []byte {
		return mockStack()
	}
	m := sync.NewRWMutex(xlog.NopLogger, "Test", 0)
	assert.Equal(t, "main.(*Worker).Retire(0xc4200f1800, 0x90cfa0, 0xc420018d70) | /home/asticode/projects/go/src/github.com/asticode/myproject/worker.go:174 +0x11d", m.Parent())
}

func TestRWMutex_Register(t *testing.T) {
	sync.Stack = func() []byte {
		return mockStack()
	}
	var i []string
	sync.Error = func(l xlog.Logger, m string, f xlog.F) {
		i = append(i, m)
	}
	m := sync.NewRWMutex(xlog.NopLogger, "Test", time.Nanosecond)
	m.Lock()
	go func() {
		time.Sleep(1 * time.Millisecond)
		m.Unlock()
	}()
	m.Lock()
	m.Unlock()
	assert.Equal(t, "Test mutex: Deadlock detected for Lock at main.(*Worker).Retire(0xc4200f1800, 0x90cfa0, 0xc420018d70) | /home/asticode/projects/go/src/github.com/asticode/myproject/worker.go:174 +0x11d", i[0])
}
