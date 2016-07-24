package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kr/pretty"
	"github.com/nimona/go-nimona/journal"
	"github.com/nimona/go-nimona/store"
)

// Aggregate
type Aggregate interface {
	GetGuid() []byte
	ApplyEvent(Event)
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	// New will return an initialized but empty concrete
	// aggregate. This is a temporary solution to reflection.
	New() Aggregate
}

type Event interface {
	GetGuid() []byte
	GetTopic() string
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
}

type EventUnknown struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

// Repository handles the events from a single Journal.
// It is responsible of managing the aggregates for the various event types.
type Repository struct {
	aggregate Aggregate
	event     Event
	pairs     map[string]Aggregate
	store     store.Store
}

func NewRepository(store store.Store, aggregate Aggregate, event Event) *Repository {
	return &Repository{
		aggregate: aggregate,
		event:     event,
		store:     store,
		pairs:     map[string]Aggregate{},
	}
}

func (r *Repository) GetByGuid(key []byte) (Aggregate, error) {
	// v, err:=r.store.Get(key)
	// if err !=nil {
	// 	return nil, err
	// }

	if a, ok := r.pairs[string(key)]; ok {
		return a, nil
	}
	return nil, errors.New("Not found") // TODO(geoah) Proper error
}

func (r *Repository) AppendedEntry(entry journal.Entry) {
	fmt.Println("> Processing", string(entry.GetIndex()), string(entry.GetPayload()))
	// TODO(geoah) Check that this event hasn't already been processed.
	// TODO(geoah) Check that the previous event has already been processed.
	aggregate := r.aggregate.New()
	event := r.event
	event.Unmarshal(entry.GetPayload())
	exists := false
	if a, ok := r.pairs[string(event.GetGuid())]; ok {
		exists = true
		aggregate = a
	}
	aggregate.ApplyEvent(event)
	if exists == false {
		r.pairs[string(event.GetGuid())] = aggregate
	}

	pretty.Println("r.pairs", r.pairs)
}
