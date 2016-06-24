package journal

import (
	"github.com/nimona/go-nimona/store"
	"github.com/nimona/go-nimona/stream"
)

// JournalClusteringKey is a composite key of 2 keys.
// The first key is the user's ID and the second key is the entry's index.
type JournalClusteringKey struct {
	keys []store.Key
}

// GetKeys returns the individual keys.
func (ck *JournalClusteringKey) GetKeys() []store.Key {
	return ck.keys
}

// IsComplete checks if the ClusteringKey is complete.
func (ck *JournalClusteringKey) IsComplete() bool {
	return len(ck.keys) == 2
}

// NewJournalClusteringKey creates a new ClusteringKey using the user's ID.
func NewJournalClusteringKey(userID store.Key, index stream.Index) store.ClusteringKey {
	return &JournalClusteringKey{
		keys: []store.Key{
			userID,
			index,
		},
	}
}

// JournalPayload is our entry payload
type JournalPayload struct {
	String string `json:"string"`
}

// JournalEntry is an Entry for the UserJournal
type JournalEntry struct {
	Index       stream.Index    `json:"index"`
	ParentIndex stream.Index    `json:"parent"`
	Payload     *JournalPayload `json:"payload"`
}

// NewJournalEntry creates a new JournalEntry
func NewJournalEntry(index, parentIndex stream.Index, payload *JournalPayload) *JournalEntry {
	return &JournalEntry{
		Index:       index,
		ParentIndex: parentIndex,
		Payload:     payload,
	}
}

// GetIndex returns the Entry's Index.
func (e *JournalEntry) GetIndex() stream.Index {
	return e.Index
}

// GetParentIndex returns the parent Entry's Index.
func (e *JournalEntry) GetParentIndex() stream.Index {
	return e.ParentIndex
}

// GetPayload returns the Payload for the Entry.
func (e *JournalEntry) GetPayload() stream.Payload {
	return e.Payload
}
