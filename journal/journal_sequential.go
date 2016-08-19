package journal

import (
	"fmt"
	"io"
	"sync"

	mc "github.com/jbenet/go-multicodec"
)

// SequentialJournal is a Journal with an incremental Index.
type SequentialJournal struct {
	sync.Mutex
	codec     mc.Codec
	notifiees []Notifiee // TODO(geoah) convert to a map so we can de-register them.
	lastIndex Index
	encoder   mc.Encoder
	decoder   mc.Decoder
}

// SequentialEntry is an Entry in our Journal with a sequential uint64 Index.
type SequentialEntry struct {
	Index   Index  `json:"i"`
	Payload []byte `json:"p"`
}

// NewSequentialEntry creates a new SequentialEntry out of an index and payload.
func NewSequentialEntry(index Index, payload []byte) *SequentialEntry {
	return &SequentialEntry{
		Index:   index,
		Payload: payload,
	}
}

// GetIndex returns the Entry's Index.
func (e *SequentialEntry) GetIndex() Index {
	return e.Index
}

// GetParentIndex returns the parent Entry's Index.
func (e *SequentialEntry) GetParentIndex() Index {
	return e.Index - 1
}

// GetPayload returns the Payload for the Entry.
func (e *SequentialEntry) GetPayload() Payload {
	return e.Payload
}

// NewJournal creates a new Journal.
func NewJournal(c mc.Codec, r io.Reader, w io.Writer) *SequentialJournal {
	dec := c.Decoder(r)
	enc := c.Encoder(w)
	return &SequentialJournal{
		lastIndex: rootEntryIndex,
		encoder:   enc,
		decoder:   dec,
	}
}

// Replay goes through all the persisted events and processes them.
func (j *SequentialJournal) Replay() error {
	for {
		e := &SequentialEntry{}
		err := j.decoder.Decode(e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		j.processEntry(false, e)
	}
	return nil
}

// Restore appends an Entry to the Journal with an existing index.
func (j *SequentialJournal) Restore(entries ...Entry) (lin Index, err error) {
	return j.processEntry(true, entries...)
}

// Append appends a payload as the next Entry to the Journal.
func (j *SequentialJournal) Append(payloads ...[]byte) (lin Index, err error) {
	for _, payload := range payloads {
		j.Lock()
		defer j.Unlock()
		entry := NewSequentialEntry(j.lastIndex+1, payload)
		if lin, err = j.processEntry(true, entry); err != nil {
			return lin, err
		}
	}
	return lin, nil
}

func (j *SequentialJournal) processEntry(persist bool, entries ...Entry) (lin Index, err error) {
	// go through all entries
	for _, entry := range entries {
		fmt.Printf("> Processing entry=%#v;\n", entry)
		// check if we have already processed the previous entry
		// or that this is the first entry of our log (index=0).
		pi := entry.GetIndex() - 1
		if pi != rootEntryIndex && pi != j.lastIndex {
			fmt.Printf("Missing parent index. lastIndex=%d;\n", j.lastIndex)
			return j.lastIndex, ErrMissingParentIndex
		}

		if persist == true {
			fmt.Printf("\t> Persisting entry.\n")
			err := j.encoder.Encode(&entry)
			if err != nil {
				fmt.Printf("\t\tCould not encode entry. err=%#v;\n", err)
				return j.lastIndex, err
			}
			fmt.Printf("\t\t Persisted entry.\n")
		}

		j.lastIndex = entry.GetIndex() // TODO(geoah) Do we need to check for type?
		j.notifyAll(entry)
	}
	return j.lastIndex, nil
}

// Notify adds notifiees for AppendEntry events.
func (j *SequentialJournal) Notify(notifiee Notifiee) {
	j.notifiees = append(j.notifiees, notifiee)
}

// notifyAll notifies anyone who cares about changes in the Journal.
func (j *SequentialJournal) notifyAll(entry Entry) {
	// TODO(geoah) Log
	fmt.Println("\t> Notifying notifiees.", entry)
	for i, notifiee := range j.notifiees {
		fmt.Printf("\t\t Notified notifiee #%d.\n", i)
		notifiee.ProcessJournalEntry(entry)
	}
}
