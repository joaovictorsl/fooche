package storage

// This is a wrapper around a map[string][]byte data structure.
type UnboundedStorage struct {
	data map[string][]byte
}

func NewUnboundedStorage() *UnboundedStorage {
	return &UnboundedStorage{
		data: make(map[string][]byte),
	}
}

func (ubs *UnboundedStorage) Put(k string, v []byte) (err error) {
	ubs.data[k] = v
	return nil
}

func (ubs *UnboundedStorage) Get(k string) (v []byte, ok bool) {
	v, ok = ubs.data[k]
	return v, ok
}

func (ubs *UnboundedStorage) Remove(k string) {
	delete(ubs.data, k)
}

func (ubs *UnboundedStorage) Size() int {
	return len(ubs.data)
}

func (ubs *UnboundedStorage) Capacity() int {
	return -1
}
