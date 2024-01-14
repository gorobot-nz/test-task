package storage

import (
	"slices"
	"sync"
)

type Storage[V any] struct {
	m sync.Mutex

	storageMap map[string]V
	keysSlice  []string
}

func NewStorage[V any]() *Storage[V] {
	return &Storage[V]{
		m:          sync.Mutex{},
		storageMap: make(map[string]V),
		keysSlice:  make([]string, 0, 1024),
	}
}

func (s *Storage[V]) Set(key string, value V) {
	s.m.Lock()
	s.storageMap[key] = value

	if !slices.Contains(s.keysSlice, key) {
		s.keysSlice = append(s.keysSlice, key)
	}
	s.m.Unlock()
}

func (s *Storage[V]) Get(key string) (V, bool) {
	val, ok := s.storageMap[key]

	return val, ok
}

func (s *Storage[V]) List() []V {
	s.m.Lock()
	var result = make([]V, len(s.keysSlice))

	for index, val := range s.keysSlice {
		result[index] = s.storageMap[val]
	}
	s.m.Unlock()

	return result
}

func (s *Storage[V]) Delete(key string) {
	s.m.Lock()
	index := slices.Index(s.keysSlice, key)
	if index != -1 {
		s.keysSlice = slices.Delete(s.keysSlice, index, index+1)
		delete(s.storageMap, key)
	}
	s.m.Unlock()
}
