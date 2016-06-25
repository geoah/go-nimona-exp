package stream

// Index is log's incremental index.
type Index interface{}

// Payload is the entry's content.
type Payload interface{}

// Entry is each of the records of our journal.
type Entry interface {
	// GetIndex returns the Entry's Index.
	GetIndex() Index
	// GetParentIndex returns the parent Entry's Index.
	GetParentIndex() Index
	// GetPayload returns the Payload for the Entry
	GetPayload() Payload
}
