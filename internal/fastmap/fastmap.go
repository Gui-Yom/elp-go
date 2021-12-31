package fastmap

import (
	"elp-go/internal/world"
	"math"
	"unsafe"
)

// intPhi is for scrambling the keys
const intPhi = 0x9E3779B9

// freeKey is the 'free' key
const freeKey = 0

func phiMix(x int64) int64 {
	h := x * intPhi
	return h ^ (h >> 16)
}

// Map is a map-like data-structure for int64s
type Map struct {
	data       []int64 // interleaved keys and values
	fillFactor float64
	threshold  int // we will resize a map once it reaches this size
	size       int

	mask  int64 // mask to calculate the original position
	mask2 int64

	hasFreeKey bool  // do we have 'free' key in the map?
	freeVal    int64 // value of 'free' key
}

func nextPowerOf2(x uint32) uint32 {
	if x == math.MaxUint32 {
		return x
	}

	if x == 0 {
		return 1
	}

	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16

	return x + 1
}

func arraySize(exp int, fill float64) int {
	s := nextPowerOf2(uint32(math.Ceil(float64(exp) / fill)))
	if s < 2 {
		s = 2
	}
	return int(s)
}

// New returns a map initialized with n spaces and uses the stated fillFactor.
// The map will grow as needed.
func New(size int, fillFactor float64) *Map {
	if fillFactor <= 0 || fillFactor >= 1 {
		panic("FillFactor must be in (0, 1)")
	}
	if size <= 0 {
		panic("Size must be positive")
	}

	capacity := arraySize(size, fillFactor)
	return &Map{
		data:       make([]int64, 2*capacity),
		fillFactor: fillFactor,
		threshold:  int(math.Floor(float64(capacity) * fillFactor)),
		mask:       int64(capacity - 1),
		mask2:      int64(2*capacity - 1),
	}
}

// Get returns the value if the key is found.
func (m *Map) get(key int64) (int64, bool) {
	if key == freeKey {
		if m.hasFreeKey {
			return m.freeVal, true
		}
		return 0, false
	}

	ptr := (phiMix(key) & m.mask) << 1
	if ptr < 0 || ptr >= int64(len(m.data)) { // Check to help to compiler to eliminate a bounds check below.
		return 0, false
	}
	k := m.data[ptr]

	if k == freeKey { // end of chain already
		return 0, false
	}
	if k == key { // we check FREE prior to this call
		return m.data[ptr+1], true
	}

	for {
		ptr = (ptr + 2) & m.mask2
		k = m.data[ptr]
		if k == freeKey {
			return 0, false
		}
		if k == key {
			return m.data[ptr+1], true
		}
	}
}

func (m *Map) GetCost(key world.Position) (float64, bool) {
	key_ := int64(key.X)<<32 | int64(key.Y)
	if val, ok := m.get(key_); ok {
		return *(*float64)(unsafe.Pointer(&val)), true
	} else {
		return 0, false
	}
}

func (m *Map) GetPos(key world.Position) (world.Position, bool) {
	key_ := int64(key.X)<<32 | int64(key.Y)
	if val, ok := m.get(key_); ok {
		X := int32(val >> 32)        // High 32bits
		Y := int32(val & 0xFFFFFFFF) // Low 32bits
		return world.Position{X: X, Y: Y}, true
	} else {
		return world.Position{}, false
	}
}

// Put adds or updates key with value val.
func (m *Map) put(key int64, val int64) {
	if key == freeKey {
		if !m.hasFreeKey {
			m.size++
		}
		m.hasFreeKey = true
		m.freeVal = val
		return
	}

	ptr := (phiMix(key) & m.mask) << 1
	k := m.data[ptr]

	if k == freeKey { // end of chain already
		m.data[ptr] = key
		m.data[ptr+1] = val
		if m.size >= m.threshold {
			m.rehash()
		} else {
			m.size++
		}
		return
	} else if k == key { // overwrite existed value
		m.data[ptr+1] = val
		return
	}

	for {
		ptr = (ptr + 2) & m.mask2
		k = m.data[ptr]

		if k == freeKey {
			m.data[ptr] = key
			m.data[ptr+1] = val
			if m.size >= m.threshold {
				m.rehash()
			} else {
				m.size++
			}
			return
		} else if k == key {
			m.data[ptr+1] = val
			return
		}
	}
}

func (m *Map) PutCost(key world.Position, cost float64) {
	key_ := int64(key.X)<<32 | int64(key.Y)
	m.put(key_, *(*int64)(unsafe.Pointer(&cost)))
}

func (m *Map) PutPos(key world.Position, pos world.Position) {
	key_ := int64(key.X)<<32 | int64(key.Y)
	val_ := int64(pos.X)<<32 | int64(pos.Y)
	m.put(key_, val_)
}

// Del deletes a key and its value.
func (m *Map) Del(key int64) {
	if key == freeKey {
		m.hasFreeKey = false
		m.size--
		return
	}

	ptr := (phiMix(key) & m.mask) << 1
	k := m.data[ptr]

	if k == key {
		m.shiftKeys(ptr)
		m.size--
		return
	} else if k == freeKey { // end of chain already
		return
	}

	for {
		ptr = (ptr + 2) & m.mask2
		k = m.data[ptr]

		if k == key {
			m.shiftKeys(ptr)
			m.size--
			return
		} else if k == freeKey {
			return
		}

	}
}

func (m *Map) shiftKeys(pos int64) int64 {
	// Shift entries with the same hash.
	var last, slot int64
	var k int64
	var data = m.data
	for {
		last = pos
		pos = (last + 2) & m.mask2
		for {
			k = data[pos]
			if k == freeKey {
				data[last] = freeKey
				return last
			}

			slot = (phiMix(k) & m.mask) << 1
			if last <= pos {
				if last >= slot || slot > pos {
					break
				}
			} else {
				if last >= slot && slot > pos {
					break
				}
			}
			pos = (pos + 2) & m.mask2
		}
		data[last] = k
		data[last+1] = data[pos+1]
	}
}

func (m *Map) rehash() {
	newCapacity := len(m.data) * 2
	m.threshold = int(math.Floor(float64(newCapacity/2) * m.fillFactor))
	m.mask = int64(newCapacity/2 - 1)
	m.mask2 = int64(newCapacity - 1)

	data := make([]int64, len(m.data)) // copy of original data
	copy(data, m.data)

	m.data = make([]int64, newCapacity)
	if m.hasFreeKey { // reset size
		m.size = 1
	} else {
		m.size = 0
	}

	var o int64
	for i := 0; i < len(data); i += 2 {
		o = data[i]
		if o != freeKey {
			m.put(o, data[i+1])
		}
	}
}

// Size returns size of the map.
func (m *Map) Size() int {
	return m.size
}

// Keys returns a channel for iterating all keys.
func (m *Map) Keys() chan int64 {
	c := make(chan int64, 10)
	go func() {
		data := m.data
		var k int64

		if m.hasFreeKey {
			c <- freeKey // value is m.freeVal
		}

		for i := 0; i < len(data); i += 2 {
			k = data[i]
			if k == freeKey {
				continue
			}
			c <- k // value is data[i+1]
		}
		close(c)
	}()
	return c
}

// Items returns a channel for iterating all key-value pairs.
func (m *Map) Items() chan [2]int64 {
	c := make(chan [2]int64, 10)
	go func() {
		data := m.data
		var k int64

		if m.hasFreeKey {
			c <- [2]int64{freeKey, m.freeVal}
		}

		for i := 0; i < len(data); i += 2 {
			k = data[i]
			if k == freeKey {
				continue
			}
			c <- [2]int64{k, data[i+1]}
		}
		close(c)
	}()
	return c
}
