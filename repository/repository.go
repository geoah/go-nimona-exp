package repository

import (
	"github.com/nimona/go-nimona/journal"
	"github.com/nimona/go-nimona/store"
)

// Resource is the value of any resource.
type Resource interface{}

// Repository handles the events from a single stream.
// It is responsible of managing the projections (resources) for the various event types.
type Repository interface {
	GetResourceByID(string) (Resource, error)
	AppendEntry(journal.Entry) (Resource, error)
	AppendedEntry(journal.Entry)
}

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
