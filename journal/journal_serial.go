package journal

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/nimona/go-nimona/store"
)

func init() {
	gob.Register(&SerialEntry{})
}

// SerialJournal is a Journal with an incremental Index.
type SerialJournal struct {
	sync.Mutex
	persistence store.Store
	notifiees   []Notifiee // TODO(geoah) convert to a map so we can de-register them.
	lastIndex   Index
	file        *os.File
	encoder     *gob.Encoder
}

// SerialEntry is an Entry in our Journal with a Index.
type SerialEntry struct {
	Index   Index
	Payload []byte
}

// NewSerialEntry creates a new SerialEntry out of an index and payload.
func NewSerialEntry(index Index, payload []byte) *SerialEntry {
	return &SerialEntry{
		Index:   index,
		Payload: payload,
	}
}

// GetIndex returns the Entry's Index.
func (e *SerialEntry) GetIndex() Index {
	return e.Index
}

// GetParentIndex returns the parent Entry's Index.
func (e *SerialEntry) GetParentIndex() Index {
	return e.Index - 1
}

// GetPayload returns the Payload for the Entry.
func (e *SerialEntry) GetPayload() Payload {
	return e.Payload
}

// NewJournal creates a new Journal.
func NewJournal(persistence store.Store) *SerialJournal {
	// f, _ := os.Create("/tmp/dat2")
	f, err := os.OpenFile("/tmp/dat2", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}

	enc := gob.NewEncoder(f)
	// defer f.Close()
	j := &SerialJournal{
		persistence: persistence,
		lastIndex:   rootEntryIndex,
		file:        f,
		encoder:     enc,
	}

	return j
}

func (j *SerialJournal) Rewind() {
	dec := gob.NewDecoder(j.file)
	for {
		var e Entry
		err := dec.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			break
		}
		// TODO(geoah) check type
		j.replay(false, e.(Entry))
	}
}

func (j *SerialJournal) getKeyForIndex(index Index) store.Key {
	return []byte(fmt.Sprintf("%d", index))
}

func (j *SerialJournal) replay(store bool, entries ...Entry) (Index, error) {
	fmt.Println(">> Replaying", entries[0].GetIndex(), string(entries[0].GetPayload()))
	entry := entries[0] // TODO(geoah) handle all entries
	pi := entry.GetIndex() - 1
	if pi != rootEntryIndex && pi != j.lastIndex {
		return j.lastIndex, ErrMissingParentIndex
	}
	// TODO(geoah) Check that entry doesn't already exist

	if store == true {
		err := j.encoder.Encode(&entry)
		if err != nil {
			fmt.Println("Could not encode gob", err)
			return j.lastIndex, err
		}
		err = j.file.Sync()
		if err != nil {
			fmt.Println("Could not flush file", err)
		}
	}

	j.lastIndex = entry.GetIndex() // TODO(geoah) Do we need to check for type?
	j.notifyAll(entry)
	return j.lastIndex, nil
}

// Restore appends an Entry to the Journal with an existing index.
func (j *SerialJournal) Restore(entries ...Entry) (Index, error) {
	return j.replay(true, entries...)
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
