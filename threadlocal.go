package threadlocal

import (
	"sync/atomic"
)

const HASH_INCREMENT = 0x61c88647

type Threadlocal struct {
	HashCode int
}

var nextHashCode = func() func() int {
	var incr = int64(HASH_INCREMENT)
	return func() int {
		return int(atomic.AddInt64(&incr, 1))
	}
}()

func New() *Threadlocal {
	return &Threadlocal{
		HashCode: nextHashCode(),
	}
}

func (tl *Threadlocal) Set(data interface{}) {
	CurrentThreadLocalMap().Set(tl, data)
}

func (tl *Threadlocal) Get() interface{} {
	return CurrentThreadLocalMap().Get(tl)
}

func (tl *Threadlocal) Remove() {
	CurrentThreadLocalMap().Remove(tl)
}
