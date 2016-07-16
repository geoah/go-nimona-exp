package journal

import (
	"sync"

	"github.com/nimona/go-nimona/store"
)

// SerialJournal is a series of entries.
type SerialJournal struct {
	sync.Mutex
	persistence store.Store
	notifiees   []Notifiee // TODO(geoah) convert to a map so we can de-register them.
	lastIndex   SerialIndex
}

// SerialIndex is serial journal's incremental index.
type SerialIndex uint64

// SerialEntry is an Entry in our Journal with a SerialIndex.
type SerialEntry struct {
	index   Index
	payload []byte
}

// NewSerialEntry creates a new SerialEntry out of an index and payload.
func NewSerialEntry(index SerialIndex, payload []byte) *SerialEntry {
	return &SerialEntry{
		index:   index,
		payload: payload,
	}
}

// GetIndex returns the Entry's Index.
func (e *SerialEntry) GetIndex() Index {
	return e.index
}

// GetParentIndex returns the parent Entry's Index.
func (e *SerialEntry) GetParentIndex() Index {
	return e.index.(SerialIndex) - 1
}

// GetPayload returns the Payload for the Entry.
func (e *SerialEntry) GetPayload() []byte {
	return e.payload
}

// NewJournal creates a new Journal.
func NewJournal(persistence store.Store) *SerialJournal {
	return &SerialJournal{
		persistence: persistence,
		lastIndex:   0,
	}
}

func (j *SerialJournal) getClusteringKeyForIndex(index Index) store.ClusteringKey {
	return NewClusteringKey(index.(SerialIndex))
}

// GetEntry returns a single Entry by it's Index.
func (j *SerialJournal) GetEntry(index Index) (Entry, error) {
	key := j.getClusteringKeyForIndex(index)
	payload, err := j.persistence.GetOne(key)
	if err != nil {
		return nil, err
	}
	entry := NewSerialEntry(index.(SerialIndex), payload)
	return entry, nil
}

// RestoreEntry appends an Entry to the Journal with an existing index.
func (j *SerialJournal) RestoreEntry(entry Entry) (Index, error) {
	if entry.GetParentIndex() != rootEntryIndex {
		_, errParent := j.GetEntry(entry.GetParentIndex())
		if errParent != nil {
			return j.lastIndex, ErrMissingParentIndex
		}
	}
	// TODO(geoah) Check that entry doesn't already exist
	key := j.getClusteringKeyForIndex(entry.GetIndex())
	errPutting := j.persistence.Put(key, entry.GetPayload())
	if errPutting != nil {
		return j.lastIndex, errPutting
	}
	j.lastIndex = entry.GetIndex().(SerialIndex) // TODO(geoah) Do we need to check for type?
	j.notifyAll(entry)
	return j.lastIndex, nil
}

// AppendEntry appends a payload as the next Entry to the Journal.
func (j *SerialJournal) AppendEntry(payload []byte) (Index, error) {
	j.Lock()
	defer j.Unlock()
	entry := NewSerialEntry(j.lastIndex+1, payload)
	return j.RestoreEntry(entry)
}

// Notify adds notifiees for AppendEntry events.
func (j *SerialJournal) Notify(notifiee Notifiee) {
	j.notifiees = append(j.notifiees, notifiee)
}

// notifyAll notifies anyone who cares about changes in the stream.
func (j *SerialJournal) notifyAll(entry Entry) {
	// TODO(geoah) Log
	// fmt.Println("> Notifying notifiees about entry", entry)
	for _, notifiee := range j.notifiees {
		notifiee.AppendedEntry(entry)
	}
}
