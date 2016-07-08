package stream

import "errors"

// ErrMissingParentIndex thrown when trying to append an Entry with its
// Parent missing.
var ErrMissingParentIndex = errors.New("Entry's parent index is missing.")

// Stream is a series of entries, serialy linked.
type Stream interface {
	// GetEntry returns a single Entry by it's Index.
	GetEntry(Index) (Entry, error)
	// AppendEntry appends an Entry to the Stream, else returns error
	// `ErrMissingParentIndex` if their parent index does not exist.
	AppendEntry(Entry) error
	// Notify registers a notifiee for signals
	Notify(Notifiee)
}

// Notifiee is an interface for an object wishing to receive
// notifications from a Network.
type Notifiee interface {
	AppendedEntry(Entry) // called when an entry has been appended
}
