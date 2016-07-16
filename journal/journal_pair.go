package journal

import (
	"github.com/nimona/go-nimona/store"
)

// ClusteringKey is a composite key of 2 keys.
// The first key is the user's ID and the second key is the entry's index.
type ClusteringKey struct {
	keys []store.Key
}

// GetKeys returns the individual keys.
func (ck *ClusteringKey) GetKeys() []store.Key {
	return ck.keys
}

// IsComplete checks if the ClusteringKey is complete.
func (ck *ClusteringKey) IsComplete() bool {
	return len(ck.keys) == 1
}

// NewClusteringKey creates a new ClusteringKey using the user's ID.
func NewClusteringKey(index SerialIndex) store.ClusteringKey {
	return &ClusteringKey{
		keys: []store.Key{
			index,
		},
	}
}
