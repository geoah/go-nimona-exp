package journal

import (
	"fmt"
	"io"
	"sync"

	mc "github.com/jbenet/go-multicodec"
)

// SerialJournal is a Journal with an incremental Index.
type SerialJournal struct {
	sync.Mutex
	codec     mc.Codec
	notifiees []Notifiee // TODO(geoah) convert to a map so we can de-register them.
	lastIndex Index
	// file      *os.File
	encoder mc.Encoder
	decoder mc.Decoder
}

// SerialEntry is an Entry in our Journal with a Index.
type SerialEntry struct {
	Index   Index  `json:"i"`
	Payload []byte `json:"p"`
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
func NewJournal(c mc.Codec, r io.Reader, w io.Writer) *SerialJournal {
	dec := c.Decoder(r)
	enc := c.Encoder(w)
	return &SerialJournal{
		lastIndex: rootEntryIndex,
		encoder:   enc,
		decoder:   dec,
	}
}

func (j *SerialJournal) Rewind() {
	for {
		e := &SerialEntry{}
		err := j.decoder.Decode(e)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			break
		}
		// TODO(geoah) check type
		j.processEntry(false, e)
	}
}

func (j *SerialJournal) processEntry(persist bool, entries ...Entry) (Index, error) {
	fmt.Println(">> processing", entries[0].GetIndex(), string(entries[0].GetPayload()))
	entry := entries[0] // TODO(geoah) handle all entries
	pi := entry.GetIndex() - 1
	if pi != rootEntryIndex && pi != j.lastIndex {
		return j.lastIndex, ErrMissingParentIndex
	}
	// TODO(geoah) Check that entry doesn't already exist

	if persist == true {
		err := j.encoder.Encode(&entry)
		if err != nil {
			fmt.Println("Could not encode entry", err)
			return j.lastIndex, err
		}
	}

	j.lastIndex = entry.GetIndex() // TODO(geoah) Do we need to check for type?
	j.notifyAll(entry)
	return j.lastIndex, nil
}

// Restore appends an Entry to the Journal with an existing index.
func (j *SerialJournal) Restore(entries ...Entry) (Index, error) {
	return j.processEntry(true, entries...)
}

// Append appends a payload as the next Entry to the Journal.
func (j *SerialJournal) Append(payloads ...[]byte) (Index, error) {
	payload := payloads[0] // TODO(geoah) handle all payloads
	j.Lock()
	defer j.Unlock()
	entry := NewSerialEntry(j.lastIndex+1, payload)
	return j.processEntry(true, entry)
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
