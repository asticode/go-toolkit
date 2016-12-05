package toolkit

import "sync"

// Canceller represents an object capable of managing a cancellation
type Canceller struct {
	cancelled bool
	channels  map[chan bool]bool
	// We're using 2 mutexes here :
	//  - using mutexCancel allows us to make sure to control when Cancel() can be called
	//  - using mutexChannels allows us to make sure concurrent accesses to canceller channels are properly done
	mutexCancel   *sync.RWMutex
	mutexChannels *sync.RWMutex
}

// NewCanceller creates a new Canceller
func NewCanceller() *Canceller {
	return &Canceller{
		channels:      make(map[chan bool]bool),
		mutexCancel:   &sync.RWMutex{},
		mutexChannels: &sync.RWMutex{},
	}
}

// Cancel cancels a process
func (c *Canceller) Cancel() {
	c.mutexCancel.Lock()
	defer c.mutexCancel.Unlock()
	c.mutexChannels.Lock()
	defer c.mutexChannels.Unlock()
	c.cancelled = true
	for ch := range c.channels {
		delete(c.channels, ch)
		close(ch)
	}
}

// Cancelled returns whether the process was cancelled
func (c *Canceller) Cancelled() bool {
	c.mutexCancel.Lock()
	defer c.mutexCancel.Unlock()
	return c.cancelled
}

// Close closes a channel
func (c *Canceller) Close(i chan bool) {
	c.mutexChannels.Lock()
	defer c.mutexChannels.Unlock()
	delete(c.channels, i)
}

// Lock locks the canceller for cancellation
func (c *Canceller) Lock() {
	c.mutexCancel.Lock()
}

// NewChannel returns a new cancellation channel
func (c *Canceller) NewChannel() (o chan bool) {
	c.mutexChannels.Lock()
	defer c.mutexChannels.Unlock()
	o = make(chan bool)
	c.channels[o] = true
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
