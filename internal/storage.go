package internal

import "sync"

type storage struct {
	data map[string]string
	sync.RWMutex
}

func NewStorage() *storage {
	return &storage{
		data: make(map[string]string),
	}
}

func (s *storage) Set(key, value string) error {
	s.Lock()
	defer s.Unlock()

	s.data[key] = value
	return nil
}

func (s *storage) Get(key string) *string {
	s.Lock()
	defer s.Unlock()

	value, ok := s.data[key]
	if !ok {
		return nil
	}

	return &value
}
