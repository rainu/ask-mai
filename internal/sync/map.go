package sync

type Map[K comparable, V any] struct {
	mutex Mutex
	_map  map[K]V
}

func (m *Map[K, V]) Get(key K) (val V) {
	m.mutex.Read(func() {
		if m._map == nil {
			return
		}
		val, _ = m._map[key]
	})
	return
}

func (m *Map[K, V]) Read(key K) (val V, ok bool) {
	m.mutex.Read(func() {
		if m._map == nil {
			ok = false
			return
		}
		val, ok = m._map[key]
	})

	return
}

func (m *Map[K, V]) Len() int {
	length := 0
	m.mutex.Read(func() {
		if m._map != nil {
			length = len(m._map)
		}
	})
	return length
}

func (m *Map[K, V]) Write(key K, val V) {
	m.mutex.Write(func() {
		if m._map == nil {
			m._map = make(map[K]V)
		}
		m._map[key] = val
	})
}

func (m *Map[K, V]) Delete(key K) {
	m.mutex.Write(func() {
		if m._map != nil {
			delete(m._map, key)
		}
	})
}

func (m *Map[K, V]) For(call func(K, V)) {
	m.mutex.Read(func() {
		if m._map == nil {
			return
		}
		for k, v := range m._map {
			call(k, v)
		}
	})
}
