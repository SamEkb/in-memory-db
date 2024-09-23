package engine

import "sync"

type Hashtable struct {
	m    sync.RWMutex
	data map[string]string
}

func NewHashtable() *Hashtable {
	return &Hashtable{data: make(map[string]string)}
}

func (h *Hashtable) Get(key string) (string, bool) {
	h.m.RLock()
	defer h.m.RUnlock()
	result, ok := h.data[key]
	return result, ok
}

func (h *Hashtable) Put(key, value string) {
	h.m.Lock()
	defer h.m.Unlock()
	h.data[key] = value
}

func (h *Hashtable) Del(key string) {
	h.m.Lock()
	defer h.m.Unlock()
	delete(h.data, key)
}
