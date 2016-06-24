package store

import (
	"encoding/json"
	"reflect"

	"sync"
)

// InMemoryStore is an implementation of an in-memory Store
type InMemoryStore struct {
	sync.Mutex
	pairs map[Key]Value
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

	s.pairs[completeKey] = string(valueJSON)
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

	for key, valueJSON := range s.pairs {
		if s.keyEqual(key, completeKey) == true {
			errUnmashal := json.Unmarshal([]byte(valueJSON.(string)), result)
			return errUnmashal
		}
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

func (s *InMemoryStore) keyEqual(key1, key2 Key) bool {
	return reflect.DeepEqual(key1, key2)
}

// NewInMemoryStore returns a new in-memory Store
func NewInMemoryStore() Store {
	return &InMemoryStore{
		pairs: map[Key]Value{},
	}
}
