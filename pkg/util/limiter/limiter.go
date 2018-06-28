package limiter

import (
	"sync"
)

// Limits the number of goroutines
// and allows you to wait
// for the worker group to complete.
type Limiter struct {
	waitGroup *sync.WaitGroup
	limit     chan struct{}
}

// Returns new limiter which restrains
// the size of worker group with number
// less or equal k.
func New(k uint) *Limiter {
	return &Limiter{
		waitGroup: &sync.WaitGroup{},
		limit:     make(chan struct{}, k),
	}
}

// Adds a new worker to groups
// and returns "done" function,
// which must be run to decrement
// worker counter.
func (lim *Limiter) Start() func() {
	lim.waitGroup.Add(1)
	lim.limit <- struct{}{}
	return func() {
		// Important!
		// Do not change order of pull
		// from channel and call of .Done()!
		<-lim.limit
		lim.waitGroup.Done()
	}
}

// Blocks until all workers complete their tasks.
func (lim *Limiter) Wait() {
	lim.waitGroup.Wait()
	close(lim.limit)
}

// Returns a number of active workers.
func (lim *Limiter) Active() int {
	return len(lim.limit)
}
