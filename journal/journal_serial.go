package journal

import (
	"fmt"
	"sync"

	"github.com/nimona/go-nimona/store"
)

// SerialJournal is a Journal with an incremental Index.
type SerialJournal struct {
	sync.Mutex
	persistence store.Store
	notifiees   []Notifiee // TODO(geoah) convert to a map so we can de-register them.
	lastIndex   Index
}

// SerialEntry is an Entry in our Journal with a Index.
type SerialEntry struct {
	index   Index
	payload []byte
}

// NewSerialEntry creates a new SerialEntry out of an index and payload.
func NewSerialEntry(index Index, payload []byte) *SerialEntry {
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
	return e.index - 1
}

// GetPayload returns the Payload for the Entry.
func (e *SerialEntry) GetPayload() Payload {
	return e.payload
}

// NewJournal creates a new Journal.
func NewJournal(persistence store.Store) *SerialJournal {
	return &SerialJournal{
		persistence: persistence,
		lastIndex:   rootEntryIndex,
	}
}

func (j *SerialJournal) getKeyForIndex(index Index) store.Key {
	return []byte(fmt.Sprintf("%d", index))
}

// Restore appends an Entry to the Journal with an existing index.
func (j *SerialJournal) Restore(entries ...Entry) (Index, error) {
	entry := entries[0] // TODO(geoah) handle all entries
	pi := entry.GetIndex() - 1
	if pi != rootEntryIndex && pi != j.lastIndex {
		return j.lastIndex, ErrMissingParentIndex
	}
	// TODO(geoah) Check that entry doesn't already exist
	key := j.getKeyForIndex(entry.GetIndex())
	errPutting := j.persistence.Put(key, []byte(entry.GetPayload()))
	if errPutting != nil {
		return j.lastIndex, errPutting
	}
	j.lastIndex = entry.GetIndex() // TODO(geoah) Do we need to check for type?
	j.notifyAll(entry)
	return j.lastIndex, nil
}

// Append appends a payload as the next Entry to the Journal.
func (j *SerialJournal) Append(payloads ...[]byte) (Index, error) {
	payload := payloads[0] // TODO(geoah) handle all payloads
	j.Lock()
	defer j.Unlock()
	entry := NewSerialEntry(j.lastIndex+1, payload)
	return j.Restore(entry)
}

// Notify adds notifiees for AppendEntry events.
func (j *SerialJournal) Notify(notifiee Notifiee) {
	j.notifiees = append(j.notifiees, notifiee)
}

// notifyAll notifies anyone who cares about changes in the Journal.
func (j *SerialJournal) notifyAll(entry Entry) {
	// TODO(geoah) Log
	// fmt.Println("> Notifying notifiees about entry", entry)
	for _, notifiee := range j.notifiees {
		notifiee.AppendedEntry(entry)
	}
}
