package engine

type Hashtable struct {
	data map[string]string
}

func NewHashtable() *Hashtable {
	return &Hashtable{data: make(map[string]string)}
}

func (h *Hashtable) Get(key string) (string, bool) {
	result, ok := h.data[key]
	return result, ok
}

func (h *Hashtable) Insert(key, value string) {
	h.data[key] = value
}

func (h *Hashtable) Del(key string) {
	delete(h.data, key)
}
