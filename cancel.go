package toolkit

import (
	"context"
	"sync"
)

// Canceller represents an object capable of managing a graceful cancellation
// We're using 2 mutexes here :
//  - using mutexCancel allows us to make sure to control when Cancel() can be called
//  - using mutexContextPool allows us to make sure concurrent accesses to context pool are properly done
type Canceller struct {
	cancelled        bool
	contextPool      map[context.Context]context.CancelFunc
	mutexCancel      *sync.RWMutex
	mutexContextPool *sync.RWMutex
}

// NewCanceller creates a new Canceller
func NewCanceller() *Canceller {
	return &Canceller{
		contextPool:      make(map[context.Context]context.CancelFunc),
		mutexCancel:      &sync.RWMutex{},
		mutexContextPool: &sync.RWMutex{},
	}
}

// Cancel cancels a process
func (c *Canceller) Cancel() {
	c.mutexCancel.Lock()
	defer c.mutexCancel.Unlock()
	c.mutexContextPool.Lock()
	defer c.mutexContextPool.Unlock()
	c.cancelled = true
	for ctx, cancel := range c.contextPool {
		c.closeUnsafe(ctx, cancel)
	}
}

// Cancelled returns whether the process was cancelled
func (c *Canceller) Cancelled() bool {
	c.mutexCancel.Lock()
	defer c.mutexCancel.Unlock()
	return c.cancelled
}

// Close closes a channel
func (c *Canceller) Close(ctx context.Context) {
	c.mutexContextPool.Lock()
	defer c.mutexContextPool.Unlock()
	if cancelFunc, ok := c.contextPool[ctx]; ok {
		c.closeUnsafe(ctx, cancelFunc)
	}
}

// closeUnsafe closes the cancel func without locking the mutex
func (c *Canceller) closeUnsafe(ctx context.Context, cancelFunc context.CancelFunc) {
	cancelFunc()
	delete(c.contextPool, ctx)
}

// Lock locks the canceller for cancellation
func (c *Canceller) Lock() {
	c.mutexCancel.Lock()
}

// NewContext returns a new context
func (c *Canceller) NewContext() (ctx context.Context) {
	c.mutexContextPool.Lock()
	defer c.mutexContextPool.Unlock()
	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithCancel(context.Background())
	c.contextPool[ctx] = cancelFunc
	return
}

// Reset resets the canceller
func (c *Canceller) Reset() {
	c.mutexCancel.Lock()
	defer c.mutexCancel.Unlock()
	c.cancelled = false
}

// Unlock unlocks the canceller for cancellation
func (c *Canceller) Unlock() {
	c.mutexCancel.Unlock()
}
