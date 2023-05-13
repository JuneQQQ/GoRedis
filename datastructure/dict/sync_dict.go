package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

const (
	EXISTS     = 0
	NOT_EXISTS = 1
	SUCCESS    = 1
)

func (s *SyncDict) Set(key string, val interface{}) (result int) {
	_, exists := s.m.Load(key)
	s.m.Store(key, val)
	if exists {
		return EXISTS
	}
	return NOT_EXISTS
}

func (s *SyncDict) Remove(key string) (result int) {
	_, exists := s.m.LoadAndDelete(key)
	if exists {
		return EXISTS
	}
	return NOT_EXISTS
}

func (s *SyncDict) Clear() (result int) {
	s.m = sync.Map{}
	return SUCCESS
}

func (s *SyncDict) Get(key string) (val interface{}, exists bool) {
	val, ok := s.m.Load(key)
	return val, ok
}

func (s *SyncDict) Len() int {
	length := 0
	s.m.Range(func(key, value any) bool {
		length++
		return true
	})
	return length
}

func (s *SyncDict) SetIfAbsent(key string, val interface{}) (result int) {
	_, exists := s.m.LoadOrStore(key, val)
	if exists {
		return EXISTS
	}
	return NOT_EXISTS
}

func (s *SyncDict) SetIfExists(key string, val interface{}) (result int) {
	_, exists := s.m.Load(key)
	if !exists {
		return NOT_EXISTS
	}
	s.m.Store(key, val)
	return EXISTS
}

func (s *SyncDict) ForEach(consumer Consumer) {
	s.m.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

func (s *SyncDict) Keys() []string {
	result := make([]string, s.Len())
	i := 0
	s.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		return true
	})
	return result
}

func (s *SyncDict) RandomKeys(count int) []string {
	result := make([]string, s.Len())
	i := 0
	s.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		if i == count {
			return false
		}
		return true
	})
	return result
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}
