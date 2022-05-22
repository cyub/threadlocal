package threadlocal

import (
	"runtime"
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
	if tl == nil {
		panic("you should not directly set the threadlocal variable to nil")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	currentThreadLocalMap().Set(tl, data)
}

func (tl *Threadlocal) Get() interface{} {
	if tl == nil {
		panic("you should not directly set the threadlocal variable to nil")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return currentThreadLocalMap().Get(tl)
}

func (tl *Threadlocal) Remove() {
	if tl == nil {
		panic("you should not directly set the threadlocal variable to nil")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	currentThreadLocalMap().Remove(tl)
}
