package reactor

import (
	"sync"
	"sync/atomic"
)

type Reactor struct {
	m     sync.RWMutex
	n     int32
	group sync.WaitGroup
	tasks []func()
}

func (reactor *Reactor) Add(task func()) {
	reactor.group.Add(1)
	atomic.AddInt32(&reactor.n, 1)
	reactor.m.Lock()
	defer reactor.m.Unlock()
	reactor.tasks = append(reactor.tasks, func() {
		defer func() {
			reactor.group.Done()
			atomic.AddInt32(&reactor.n, -1)
		}()
		task()
	})
}

func (reactor *Reactor) Run() {
	func() {
		reactor.m.RLock()
		defer reactor.m.RUnlock()
		for _, task := range reactor.tasks {
			go task()
		}
		reactor.tasks = []func(){}
	}()
	reactor.group.Wait()
}

func (reactor *Reactor) Done() bool {
	return reactor.n == 0
}

func (reactor *Reactor) Tasks() int {
	return int(reactor.n)
}
