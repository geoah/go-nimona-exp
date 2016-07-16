package store

import "sync"

// InMemoryStore is an implementation of an in-memory Store
type InMemoryStore struct {
	sync.Mutex
	pairs map[string]Value
}

func (s *InMemoryStore) key(key Key) string {
	return string(key)
}

// Put sets the key's value, overwriting the previous if it exists.
func (s *InMemoryStore) Put(key Key, value Value) (err error) {
	s.Lock()
	defer s.Unlock()

	s.pairs[s.key(key)] = value
	return nil
}

// GetOne gets the value for a clustering key and updates theresult, else
// errors with `ErrKeyNotFound`.
func (s *InMemoryStore) GetOne(key Key) (value Value, err error) {
	s.Lock()
	defer s.Unlock()

	if value, ok := s.pairs[s.key(key)]; ok {
		return value, nil
	}

	// result = value
	return nil, ErrKeyNotFound
}

// GetAll finds all pairs that partially match the given key (left to right).
func (s *InMemoryStore) GetAll(key Key) (results []*Value, err error) {
	// TODO Implement
	return results, nil
}

// Delete removed the key's value if it exists, else errors with
// `ErrKeyNotFound`, or `ErrKeyNotComplete`.
func (s *InMemoryStore) Delete(key Key) (err error) {
	// TODO Implement
	return nil
}

// NewInMemoryStore returns a new in-memory Store
func NewInMemoryStore() Store {
	return &InMemoryStore{
		pairs: map[string]Value{},
	}
}
