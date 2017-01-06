package sync

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"flag"

	"bytes"
	"regexp"

	"github.com/rs/xid"
	"github.com/rs/xlog"
)

// Constants
const (
	loggerKeyMutexID = "mutex_id"
)

// Flags
var (
	byteLineDelimiter    = []byte("\n")
	mutexDeadlockTimeout = flag.Duration("mutex-deadlock-timeout", 0, "the mutex deadlock timeout")
	regexpStack          = regexp.MustCompile("github\\.com\\/asticode\\/go-sync")
)

// RWMutexLogged represents a RWMutex capable of logging its actions to ease deadlock debugging
type RWMutexLogged struct {
	Logger xlog.Logger
	Name   string
	*sync.RWMutex
	Timeout time.Duration
}

// RWMutexInfo represents the info of a RWMutex
type RWMutexInfo struct {
	Action string
	C      chan bool
	ID     string
	Parent string
}

// NewRWMutexLogged creates a new RWMutexLogged
func NewRWMutexLogged(l xlog.Logger, name string, t time.Duration) *RWMutexLogged {
	return &RWMutexLogged{
		Logger:  l,
		Name:    name,
		RWMutex: &sync.RWMutex{},
		Timeout: t,
	}
}

// NewRWMutexLoggedFromFlag creates a new RWMutexLogged based on the flag config
func NewRWMutexLoggedFromFlag(l xlog.Logger, name string) *RWMutexLogged {
	return NewRWMutexLogged(l, name, *mutexDeadlockTimeout)
}

// Stack allows testing functions using it
var Stack = func() []byte {
	return debug.Stack()
}

// Error allows testing functions using it
var Error = func(l xlog.Logger, m string, f xlog.F) {
	l.Error(m, f)
}

// Debug allows testing functions using it
var Debug = func(l xlog.Logger, m string, f xlog.F) {
	l.Debug(m, f)
}

// Parent returns the parent of the mutex
func (m *RWMutexLogged) Parent() (o string) {
	s := bytes.Split(Stack(), byteLineDelimiter)
	if len(s) > 3 {
		for a := 3; a < len(s); a++ {
			if len(regexpStack.Find(s[a])) == 0 {
				o = strings.Trim(string(s[a]), "\t")
				if len(s) >= a+2 {
					o += " | " + strings.Trim(string(s[a+1]), "\t")
				}
				return
			}
		}
	}
	return
}

// Log logs mutex related information
func (m *RWMutexLogged) Register(action string) (o RWMutexInfo) {
	// Init
	o = RWMutexInfo{
		Action: action,
		C:      make(chan bool),
		ID:     xid.New().String(),
		Parent: m.Parent(),
	}

	// Log
	Debug(m.Logger, fmt.Sprintf("%s mutex: %s requested at %s", m.Name, action, o.Parent), xlog.F{loggerKeyMutexID: o.ID})

	// Spawn go routine to detect deadlock
	go func(r RWMutexInfo, timeout time.Duration) {
		var timedOut bool
		for {
			select {
			case <-r.C:
				Debug(m.Logger, fmt.Sprintf("%s mutex: %s delivered at %s", m.Name, r.Action, r.Parent), xlog.F{loggerKeyMutexID: r.ID})
				return
			case <-time.After(timeout):
				if !timedOut {
					Error(m.Logger, fmt.Sprintf("%s mutex: Deadlock detected for %s at %s", m.Name, r.Action, r.Parent), xlog.F{loggerKeyMutexID: r.ID})
				}
				timedOut = true
			}
		}
	}(o, m.Timeout)
	return
}

// Lock write locks the mutex
func (m *RWMutexLogged) Lock() {
	if m.Timeout > 0 {
		mi := m.Register("Lock")
		defer func() {
			mi.C <- true
		}()
	}
	m.RWMutex.Lock()
	return
}

// Unlock write unlocks the mutex
func (m *RWMutexLogged) Unlock() {
	m.RWMutex.Unlock()
}

// RLock read locks the mutex
func (m *RWMutexLogged) RLock() {
	if m.Timeout > 0 {
		mi := m.Register("RLock")
		defer func() {
			mi.C <- true
		}()
	}
	m.RWMutex.Lock()
	return
}

// RUnlock read unlocks the mutex
func (m *RWMutexLogged) RUnlock() {
	m.RWMutex.Unlock()
}
