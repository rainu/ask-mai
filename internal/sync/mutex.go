package sync

import "sync"

type Mutex struct {
	sync.RWMutex
}

func (m *Mutex) Read(process func()) {
	m.RLock()
	defer m.RUnlock()

	process()
}

func (m *Mutex) Write(process func()) {
	m.Lock()
	defer m.Unlock()

	process()
}
