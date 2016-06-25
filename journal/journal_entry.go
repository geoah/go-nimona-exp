package journal

import "github.com/nimona/go-nimona/stream"

// Entry is an Entry for the UserJournal
type Entry struct {
	Index   Index    `json:"index"`
	Payload *Payload `json:"payload"`
}

// NewEntry creates a new Entry
func NewEntry(index Index, payload *Payload) *Entry {
	return &Entry{
		Index:   index,
		Payload: payload,
	}
}

// GetIndex returns the Entry's Index.
func (e *Entry) GetIndex() stream.Index {
	return e.Index
}

// GetParentIndex returns the parent Entry's Index.
func (e *Entry) GetParentIndex() stream.Index {
	return e.Index - 1
}

// GetPayload returns the Payload for the Entry.
func (e *Entry) GetPayload() stream.Payload {
	return e.Payload
}
