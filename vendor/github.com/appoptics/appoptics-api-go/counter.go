package appoptics

import "sync/atomic"

// SynchronizedCounter wraps an int64 with functions to concurrently add/increment
// and to periodically query & reset the sum.
type SynchronizedCounter int64

// NewCounter returns a new SynchronizedCounter initialized to 0.
func NewCounter() *SynchronizedCounter {
	c := SynchronizedCounter(0)
	return &c
}

// Incr adds 1 to the counter.
func (c *SynchronizedCounter) Incr() {
	c.Add(1)
}

// Add adds the specified delta to the counter.
func (c *SynchronizedCounter) Add(delta int64) {
	atomic.AddInt64((*int64)(c), delta)
}

// AddInt is a convenience function to add delta to the counter, where delta is an int.
func (c *SynchronizedCounter) AddInt(delta int) {
	c.Add(int64(delta))
}

// Reset returns the current value and resets the counter to zero
func (c *SynchronizedCounter) Reset() int64 {
	return atomic.SwapInt64((*int64)(c), 0)
}
