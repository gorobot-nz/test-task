package storage

import "slices"

type Storage[V any] struct {
	storageMap map[string]V
	keysSlice  []string
}

func NewStorage[V any]() *Storage[V] {
	return &Storage[V]{
		storageMap: make(map[string]V),
		keysSlice:  make([]string, 0, 1024),
	}
}

func (s *Storage[V]) Set(key string, value V) {
	s.storageMap[key] = value

	if !slices.Contains(s.keysSlice, key) {
		s.keysSlice = append(s.keysSlice, key)
	}
}

func (s *Storage[V]) Get(key string) V {
	return s.storageMap[key]
}

func (s *Storage[V]) List() []V {
	var result = make([]V, len(s.keysSlice))

	for index, val := range s.keysSlice {
		result[index] = s.storageMap[val]
	}

	return result
}

func (s *Storage[V]) Delete(key string) {
	index := slices.Index(s.keysSlice, key)
	if index != -1 {
		slices.Delete(s.keysSlice, index, index+1)
		delete(s.storageMap, key)
	}
}
