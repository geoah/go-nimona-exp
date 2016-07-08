package journal

import (
	"github.com/nimona/go-nimona/store"
	"github.com/nimona/go-nimona/stream"
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
	return len(ck.keys) == 2
}

// NewClusteringKey creates a new ClusteringKey using the user's ID.
func NewClusteringKey(userID store.Key, index stream.Index) store.ClusteringKey {
	return &ClusteringKey{
		keys: []store.Key{
			userID,
			index,
		},
	}
}

// Payload is our Journal's entry payload
type Payload []byte

// Index is our Journal's entry index
type Index uint64
