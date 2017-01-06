package sync

import (
	"sync"
	"time"

	"github.com/asticode/go-toolkit/debug"
	"github.com/rs/xlog"
)

// Constants
const (
	loggerKeyMutexName = "mutex_name"
)

// RWMutex represents a RWMutex capable of logging its actions to ease deadlock debugging
type RWMutex struct {
	lastSuccessfulLockStack debug.Stack
	logger                  xlog.Logger
	mutex                   *sync.RWMutex
	name                    string
}

// NewRWMutex creates a new RWMutex
func NewRWMutex(l xlog.Logger, name string) *RWMutex {
	return &RWMutex{
		logger: l,
		mutex:  &sync.RWMutex{},
		name:   name,
	}
}

// Lock write locks the mutex
func (m *RWMutex) Lock() {
	m.logger.Debugf("Requesting lock for %s", m.name, xlog.F{
		loggerKeyMutexName: m.name,
	})
	m.mutex.Lock()
	m.logger.Debugf("Lock acquired for %s", m.name, xlog.F{
		loggerKeyMutexName: m.name,
	})
	m.lastSuccessfulLockStack = debug.NewStack()
}

// Unlock write unlocks the mutex
func (m *RWMutex) Unlock() {
	m.mutex.Unlock()
	m.logger.Debugf("Unlock executed for %s", m.name, xlog.F{
		loggerKeyMutexName: m.name,
	})
}

// RLock read locks the mutex
func (m *RWMutex) RLock() {
	m.logger.Debugf("Requesting rlock for %s", m.name, xlog.F{
		loggerKeyMutexName: m.name,
	})
	m.mutex.RLock()
	m.logger.Debugf("RLock acquired for %s", m.name, xlog.F{
		loggerKeyMutexName: m.name,
	})
	m.lastSuccessfulLockStack = debug.NewStack()
}

// RUnlock read unlocks the mutex
func (m *RWMutex) RUnlock() {
	m.mutex.Unlock()
	m.logger.Debugf("RUnlock executed for %s", m.name, xlog.F{
		loggerKeyMutexName: m.name,
	})
}

// IsDeadlocked checks whether the mutex is deadlocked with a given timeout
func (m *RWMutex) IsDeadlocked(timeout time.Duration) (o bool) {
	o = true
	var channelLockAcquired = make(chan bool)
	go func() {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		channelLockAcquired <- true
	}()
	for {
		select {
		case <-channelLockAcquired:
			o = false
			return
		case <-time.After(timeout):
			return
		}
	}
	return
}

// LastSuccessfulLockCaller returns the stack item of the last successful lock caller
func (m *RWMutex) LastSuccessfulLockCaller() (s debug.StackItem) {
	if len(m.lastSuccessfulLockStack) >= 4 {
		s = m.lastSuccessfulLockStack[4]
	}
	return
}
