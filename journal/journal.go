package journal

import "errors"

// ErrMissingParentIndex thrown when trying to append an Entry with its
// Parent missing.
var ErrMissingParentIndex = errors.New("Entry's parent index is missing.")

// Index is a unique identifier of each entry.
type Index interface{}

// Payload is the value of each entry.
type Payload []byte

// Entry is each of the records of our journal.
type Entry interface {
	// GetIndex returns the Entry's Index.
	GetIndex() Index
	// GetParentIndex returns the Entry's previous/parent Index.
	// This is mainly useful when Entries do not have a sequential
	// Index. eg. In case of a DAG based Journal.
	GetParentIndex() Index
	// GetPayload returns the Payload for the Entry
	GetPayload() Payload
}

// Notifiee is an interface for an object wishing to receive notifications
// of appended Entries in a Journal.
type Notifiee interface {
	// ProcessJournalEntry will be called when an entry has been appended
	// and persisted in the journal.
	ProcessJournalEntry(Entry)
}

// Journal is a series of entries.
type Journal interface {
	// Append appends a Payload to the Journal as the next Index,
	// it returns the new index under which the Entry was added,
	// or returns `ErrMissingParentIndex` if their parent Index does not exist.
	// Append will notify all Notifiees about the appended Entry.
	Append(payload ...Payload) (Index, error)
	// Restore restores an existing Entry to the Journal,
	// or returns `ErrMissingParentIndex` if their parent index does not exist.
	// Restore is used when replicating a Journal and will return the last known
	// Index when returning `ErrMissingParentIndex` so the sender can replay
	// all Entries since that Index.
	// Restore will notify all Notifiees about the restored Entry.
	Restore(entry ...Entry) (Index, error)
	// Notify registers a notifiee that is interested in when new Entries have
	// been appended in the Journal.
	// Notify can specify the index from which notifications will start in order
	// to skip Entries that might have already been processed.
	Notify(Notifiee, Index)
}
