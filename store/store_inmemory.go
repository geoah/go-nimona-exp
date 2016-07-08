package store

import (
	"encoding/json"

	"sync"
)

// InMemoryStore is an implementation of an in-memory Store
type InMemoryStore struct {
	sync.Mutex
	pairs map[string]Value
}

func (s *InMemoryStore) key(completeKey ClusteringKey) string {
	key, _ := json.Marshal(completeKey.GetKeys())
	return string(key)
}

// Put sets the key's value, overwriting the previous if it exists.
func (s *InMemoryStore) Put(completeKey ClusteringKey, value Value) (err error) {
	if completeKey.IsComplete() == false {
		return ErrClusteringKeyNotComplete
	}

	s.Lock()
	defer s.Unlock()

	valueJSON, errJSON := json.Marshal(value)
	if errJSON != nil {
		return ErrInternalError // TODO(geoah) ErrMarshalling
	}

	s.pairs[s.key(completeKey)] = string(valueJSON)
	return nil
}

// GetOne gets the value for a clustering key and updates theresult, else
// errors with `ErrClusteringKeyNotFound`, or `ErrClusteringKeyNotComplete`.
func (s *InMemoryStore) GetOne(completeKey ClusteringKey, result Value) (err error) {
	if completeKey.IsComplete() == false {
		return ErrClusteringKeyNotComplete
	}

	s.Lock()
	defer s.Unlock()

	if valueJSON, ok := s.pairs[s.key(completeKey)]; ok {
		errUnmashal := json.Unmarshal([]byte(valueJSON.(string)), result)
		return errUnmashal
	}

	// result = value
	return ErrClusteringKeyNotFound
}

// GetAll updates the results list with the values of the given incomplete
// ClusteringKey, or `ErrClusteringKeyComplete`
func (s *InMemoryStore) GetAll(key ClusteringKey, results []*Value) (err error) {
	// TODO Implement
	return nil
}

// Delete removed the key's value if it exists, else errors with
// `ErrClusteringKeyNotFound`, or `ErrClusteringKeyNotComplete`.
func (s *InMemoryStore) Delete(completeKey ClusteringKey) (err error) {
	// TODO Implement
	return nil
}

// NewInMemoryStore returns a new in-memory Store
func NewInMemoryStore() Store {
	return &InMemoryStore{
		pairs: map[string]Value{},
	}
}
