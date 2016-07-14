package journal

import (
	"errors"
	"fmt"

	"github.com/nimona/go-nimona/store"
)

const rootEntryIndex Index = 0

// ErrMissingParentIndex thrown when trying to append an Entry with its
// Parent missing.
var ErrMissingParentIndex = errors.New("Entry's parent index is missing.")

// Journal is a series of entries.
// type Journal interface {
// 	// GetEntry returns a single Entry by it's Index.
// 	GetEntry(Index) (Entry, error)
// 	// AppendEntry appends an Entry to the Stream, else returns error
// 	// `ErrMissingParentIndex` if their parent index does not exist.
// 	AppendEntry(Entry) error
// 	// Notify registers a notifiee for signals
// 	Notify(Notifiee)
// }

// Notifiee is an interface for an object wishing to receive
// notifications from a Network.
type Notifiee interface {
	AppendedEntry(Entry) // called when an entry has been appended
}

// Journal is a series of entries.
type Journal struct {
	persistence store.Store
	notifiees   []Notifiee // TODO(geoah) convert to a map so we can de-register them.
	lastIndex   Index
}

// NewJournal creates a new Journal.
func NewJournal(persistence store.Store) *Journal {
	return &Journal{
		persistence: persistence,
		lastIndex:   0,
	}
}

func (j *Journal) getClusteringKeyForIndex(index Index) store.ClusteringKey {
	return NewClusteringKey(index)
}

// GetEntry returns a single Entry by it's Index.
func (j *Journal) GetEntry(index Index) (Entry, error) {
	key := j.getClusteringKeyForIndex(index)
	entry := &BasicEntry{}
	err := j.persistence.GetOne(key, entry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// AppendEntry appends an Entry to the Journal.
func (j *Journal) AppendEntry(entry Entry) error {
	if entry.GetParentIndex() != rootEntryIndex {
		_, errParent := j.GetEntry(entry.GetParentIndex())
		if errParent != nil {
			return ErrMissingParentIndex
		}
	}
	// TODO(geoah) Check that entry doesn't already exist
	key := j.getClusteringKeyForIndex(entry.GetIndex())
	errPutting := j.persistence.Put(key, entry)
	if errPutting != nil {
		return errPutting
	}
	j.lastIndex = entry.GetIndex()
	j.notifyAll(entry)
	return nil
}

// AppendPayload appends a payload as the Entry to the Journal.
func (j *Journal) AppendPayload(payload Payload) error {
	entry := NewEntry(j.lastIndex+1, payload)
	return j.AppendEntry(entry)
}

// Notify adds notifiees for AppendEntry events.
func (j *Journal) Notify(notifiee Notifiee) {
	j.notifiees = append(j.notifiees, notifiee)
}

// Notify notifies anyone who cares about changes in the stream.
func (j *Journal) notifyAll(entry Entry) {
	fmt.Println("> Notifying notifiees about entry", entry)
	for _, notifiee := range j.notifiees {
		notifiee.AppendedEntry(entry)
	}
}
