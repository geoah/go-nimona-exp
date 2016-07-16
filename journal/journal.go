package journal

import (
	"errors"

	"github.com/nimona/go-nimona/store"
)

const rootEntryIndex SerialIndex = 0

// ErrMissingParentIndex thrown when trying to append an Entry with its
// Parent missing.
var ErrMissingParentIndex = errors.New("Entry's parent index is missing.")

// Journal is a series of entries.
type Journal interface {
	// GetEntry returns a single Entry by it's Index.
	GetEntry(Index) (Entry, error)
	// AppendEntry appends an Entry to the Stream, else returns error
	// `ErrMissingParentIndex` if their parent index does not exist.
	AppendEntry(Entry) error
	// Notify registers a notifiee for signals
	Notify(Notifiee)
}

// Index is the journal's index.
type Index interface{}

// Entry is each of the records of our journal.
type Entry interface {
	// GetIndex returns the Entry's Index.
	GetIndex() Index
	// GetParentIndex returns the parent Entry's Index.
	GetParentIndex() Index
	// GetPayload returns the Payload for the Entry
	GetPayload() []byte
}

// Notifiee is an interface for an object wishing to receive
// notifications from a Network.
type Notifiee interface {
	AppendedEntry(Entry) // called when an entry has been appended
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

// NewClusteringKey creates a new ClusteringKey using the user's ID.
func NewClusteringKey(index SerialIndex) store.ClusteringKey {
	return &ClusteringKey{
		keys: []store.Key{
			index,
		},
	}
}
