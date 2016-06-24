package journal

import (
	"github.com/nimona/go-nimona/store"
	"github.com/nimona/go-nimona/stream"
)

// Journal is a Stream for user journals.
type Journal struct {
	userID      store.Key
	persistence store.Store
}

// NewJournal creates a new Journal.
func NewJournal(userID store.Key, persistence store.Store) *Journal {
	rootEntry := NewJournalEntry("0", "0", &JournalPayload{})
	key := NewJournalClusteringKey(userID, rootEntry.GetIndex())
	persistence.Put(key, rootEntry)
	return &Journal{
		userID:      userID,
		persistence: persistence,
	}
}

func (j *Journal) getClusteringKeyForIndex(index stream.Index) store.ClusteringKey {
	return NewJournalClusteringKey(j.userID, index)
}

// GetEntry returns a single Entry by it's Index.
func (j *Journal) GetEntry(index stream.Index) (stream.Entry, error) {
	key := j.getClusteringKeyForIndex(index)
	entry := &JournalEntry{}
	err := j.persistence.GetOne(key, entry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// AppendEntry appends an Entry to the Journal.
func (j *Journal) AppendEntry(entry stream.Entry) error {
	_, errParent := j.GetEntry(entry.GetParentIndex())
	if errParent != nil {
		return stream.ErrMissingParentIndex
	}
	key := j.getClusteringKeyForIndex(entry.GetIndex())
	errPutting := j.persistence.Put(key, entry)
	if errPutting != nil {
		return errPutting
	}
	return nil
}
