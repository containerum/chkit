package reactor

import (
	"sync"
	"sync/atomic"
)

type Reactor struct {
	n     int32
	group sync.WaitGroup
	tasks []func()
}

func (reactor *Reactor) Add(task func()) {
	reactor.group.Add(1)
	reactor.n++
	reactor.tasks = append(reactor.tasks, func() {
		defer func() {
			reactor.group.Done()
			atomic.AddInt32(&reactor.n, -1)
		}()
		task()
	})
}

func (reactor *Reactor) Run() {
	go func() {
		for _, task := range reactor.tasks {
			go task()
		}
	}()
	reactor.group.Wait()
}
