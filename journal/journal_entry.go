package journal

// SerialIndex is serial journal's incremental index.
type SerialIndex uint64

// Payload is the entry's content.
// type Payload []byte

// Entry is each of the records of our journal.
type Entry interface {
	// GetIndex returns the Entry's Index.
	GetIndex() Index
	// GetParentIndex returns the parent Entry's Index.
	GetParentIndex() Index
	// GetPayload returns the Payload for the Entry
	GetPayload() []byte
}

// SerialEntry is an Entry in our Journal with a SerialIndex.
type SerialEntry struct {
	index   Index
	payload []byte
}

// NewSerialEntry creates a new SerialEntry out of an index and payload.
func NewSerialEntry(index SerialIndex, payload []byte) *SerialEntry {
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
	return e.index.(SerialIndex) - 1
}

// GetPayload returns the Payload for the Entry.
func (e *SerialEntry) GetPayload() []byte {
	return e.payload
}
