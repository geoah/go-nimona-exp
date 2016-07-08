package journal

import (
	"encoding/json"
	"fmt"

	"github.com/nimona/go-nimona/store"
	"github.com/nimona/go-nimona/stream"
)

var rootEntryIndex Index

// Journal is a Stream for user journals.
type Journal struct {
	userID      store.Key
	persistence store.Store
	notifiees   []stream.Notifiee // TODO(geoah) convert to a map so we can de-register them.
	lastIndex   Index
}

// NewJournal creates a new Journal.
func NewJournal(userID store.Key, persistence store.Store) *Journal {
	return &Journal{
		userID:      userID,
		persistence: persistence,
		lastIndex:   0,
	}
}

func (j *Journal) getClusteringKeyForIndex(index stream.Index) store.ClusteringKey {
	return NewClusteringKey(j.userID, index)
}

// GetEntry returns a single Entry by it's Index.
func (j *Journal) GetEntry(index stream.Index) (stream.Entry, error) {
	key := j.getClusteringKeyForIndex(index)
	entry := &Entry{}
	err := j.persistence.GetOne(key, entry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// AppendEntry appends an Entry to the Journal.
func (j *Journal) AppendEntry(entry stream.Entry) error {
	if entry.GetParentIndex() != rootEntryIndex {
		_, errParent := j.GetEntry(entry.GetParentIndex())
		if errParent != nil {
			return stream.ErrMissingParentIndex
		}
	}
	// TODO(geoah) Check that entry doesn't already exist
	key := j.getClusteringKeyForIndex(entry.GetIndex())
	errPutting := j.persistence.Put(key, entry)
	if errPutting != nil {
		return errPutting
	}
	j.lastIndex = entry.GetIndex().(Index)
	j.notifyAll(entry)
	return nil
}

// AppendPayload appends a payload as the Entry to the Journal.
func (j *Journal) AppendPayload(payload stream.Payload) error {
	// TODO(geoah) Why do we even accept anything other than []byte?
	payloadBytes := []byte{}
	if pb, ok := payload.([]byte); ok {
		payloadBytes = pb
	} else {
		var err error
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}
	entry := NewEntry(j.lastIndex+1, payloadBytes)
	return j.AppendEntry(entry)
}

// Notify adds notifiees for AppendEntry events.
func (j *Journal) Notify(notifiee stream.Notifiee) {
	j.notifiees = append(j.notifiees, notifiee)
}

// Notify notifies anyone who cares about changes in the stream.
func (j *Journal) notifyAll(entry stream.Entry) {
	fmt.Println("> Notifying notifiees about entry", entry)
	for _, notifiee := range j.notifiees {
		notifiee.AppendedEntry(entry)
	}
}
