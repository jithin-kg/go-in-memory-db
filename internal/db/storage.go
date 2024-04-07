package db

import "sync"

type KeyValueStore struct {
	store map[string]interface{} // for read write lock
	mu    sync.Mutex
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		store: make(map[string]interface{}),
	}
}

func (s *KeyValueStore) Set(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = value
	return nil
}

func (s *KeyValueStore) Get(key string) (value interface{}, exists bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists = s.store[key]
	return
}

// delete an entry from map based on the key
func (s *KeyValueStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, key)
}
