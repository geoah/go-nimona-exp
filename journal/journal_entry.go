package journal

// Index is log's incremental index.
type Index uint64

// Payload is the entry's content.
type Payload []byte

// Entry is each of the records of our journal.
type Entry interface {
	// GetIndex returns the Entry's Index.
	GetIndex() Index
	// GetParentIndex returns the parent Entry's Index.
	GetParentIndex() Index
	// GetPayload returns the Payload for the Entry
	GetPayload() Payload
}

// BasicEntry is an Entry for the UserJournal
type BasicEntry struct {
	Index   Index  `json:"index"`
	Payload []byte `json:"payload"`
}

// NewEntry creates a new Entry
func NewEntry(index Index, payload []byte) *BasicEntry {
	return &BasicEntry{
		Index:   index,
		Payload: payload,
	}
}

// GetIndex returns the Entry's Index.
func (e *BasicEntry) GetIndex() Index {
	return e.Index
}

// GetParentIndex returns the parent Entry's Index.
func (e *BasicEntry) GetParentIndex() Index {
	return e.Index - 1
}

// GetPayload returns the Payload for the Entry.
func (e *BasicEntry) GetPayload() Payload {
	return e.Payload
}
