package atomsk

import "sync/atomic"

type Bool struct {
	value atomic.Value
}

func (b *Bool) Store(v bool) *Bool {
	b.value.Store(v)
	return b
}

func (b *Bool) Bool() bool {
	return b.value.Load().(bool)
}

func (b *Bool) True() bool {
	return b.Bool()
}

func (b *Bool) False() bool {
	return !b.Bool()
}
