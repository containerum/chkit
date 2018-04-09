package animation

import (
	"runtime"
	"sync/atomic"
)

type futuristicStop struct {
	isInitialised uint64
	isStopped     uint64
	stop          chan struct{}
}

func (fstop *futuristicStop) Init() func() {
	fstop.stop = make(chan struct{})
	atomic.CompareAndSwapUint64(&fstop.isInitialised, 0, 1)
	return func() {
		atomic.CompareAndSwapUint64(&fstop.isStopped, 0, 1)
	}
}

func (fstop *futuristicStop) Stop() {
	for !atomic.CompareAndSwapUint64(&fstop.isInitialised, 1, 1) {
		runtime.Gosched()
	}
	close(fstop.stop)
	for !atomic.CompareAndSwapUint64(&fstop.isStopped, 1, 1) {
		runtime.Gosched()
	}
}

func (fstop *futuristicStop) Done() <-chan struct{} {
	return fstop.stop
}
