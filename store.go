package threadlocal

const (
	INITIALIZE_THREADLOCALMAP_SIZE   = 16
	INITIALIZE_THREADLOCALSTORE_SIZE = 128
)

type threadLocalStore map[uint32]*ThreadlocalMap

var (
	store    = make(threadLocalStore, INITIALIZE_THREADLOCALSTORE_SIZE)
	storeMux mutex
)

func NewThreadlocalMap(capacity int) *ThreadlocalMap {
	return &ThreadlocalMap{
		size:     0,
		capacity: capacity,
		entities: make([]*Entity, capacity, capacity),
	}
}

type ThreadlocalMap struct {
	size     int
	capacity int
	entities []*Entity
}

func (tlm *ThreadlocalMap) Size() int {
	return tlm.size
}

func (tlm *ThreadlocalMap) Set(key *Threadlocal, val interface{}) {
	i := key.HashCode & (tlm.capacity - 1)
	for e := tlm.entities[i]; e != nil; i = tlm.nextIndex(i, tlm.capacity) {
		if e.key == key {
			e.val = val
			return
		}
	}
	tlm.entities[i] = NewEntity(key, val)
	tlm.size++
	tlm.rehash()
}

func (tlm *ThreadlocalMap) Get(key *Threadlocal) interface{} {
	i := key.HashCode & (tlm.capacity - 1)
	for e := tlm.entities[i]; e != nil; i = tlm.nextIndex(i, tlm.capacity) {
		if e.key == key {
			return e.val
			break
		}
	}
	return nil
}

func (tlm *ThreadlocalMap) Remove(key *Threadlocal) {
	i := key.HashCode & (tlm.capacity - 1)
	for e := tlm.entities[i]; e != nil; i = tlm.nextIndex(i, tlm.capacity) {
		if e.key == key {
			tlm.entities[i] = nil
			tlm.moveNextSameSlotEntity((key.HashCode & (tlm.capacity - 1)), i)
		}
	}
}

func (tlm *ThreadlocalMap) moveNextSameSlotEntity(slot int, idx int) {
	current := idx
	i := idx + 1
	for e := tlm.entities[i]; e != nil; i = tlm.nextIndex(i, tlm.capacity) {
		if e.key.HashCode&(tlm.capacity-1) != slot {
			break
		}
		tlm.entities[current] = e
		current = i
		tlm.entities[current] = nil
	}
}

func (tlm *ThreadlocalMap) nextIndex(i, cap int) int {
	if i+1 < cap {
		return i + 1
	}
	return 0
}

func (tlm *ThreadlocalMap) rehash() {
	if tlm.size < tlm.capacity*3/4 {
		return
	}
	newLen := tlm.capacity << 1
	newTab := make([]*Entity, newLen, newLen)
	for _, e := range tlm.entities {
		if e == nil {
			continue
		}
		i := e.key.HashCode & (newLen - 1)
		for newTab[i] != nil {
			i = tlm.nextIndex(i, newLen)
		}
		newTab[i] = e
	}
	tlm.capacity = newLen
	tlm.entities = newTab
}

type Entity struct {
	key *Threadlocal
	val interface{}
}

func NewEntity(tl *Threadlocal, val interface{}) *Entity {
	return &Entity{
		key: tl,
		val: val,
	}
}

func CurrentThreadLocalMap() *ThreadlocalMap {
	var tid uint32
	tid = ThreadId()
	lock(&storeMux)
	if store[tid] == nil {
		store[tid] = NewThreadlocalMap(INITIALIZE_THREADLOCALMAP_SIZE)
	}
	unlock(&storeMux)
	return store[tid]
}
