package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/kr/pretty"
	"github.com/nimona/go-nimona/journal"
	"github.com/nimona/go-nimona/store"
)

var (
	// ErrAggregateNotFound is returned when a requested `Aggregate`
	// cannot be found.
	ErrAggregateNotFound = errors.New("Not found")
)

// Aggregate is the result of applying a series of events.
// eg. A `User` could be the `Aggregate` of `UserCreated`,
// `UserChangedPassword`, and `UserChangedEmail` events.
type Aggregate interface {
	GetGuid() []byte
	ApplyEvent(Event)
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
}

// Event is anything that has happened.
// eg. `UserCreated`, `UserChangedPassword`, `UserChangedEmail`
type Event interface {
	GetGuid() []byte
	GetTopic() string
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
}

// EventUnknown is used internally to allow us to retrieve the
// event's topic without knowing or caring what its payload is.
// TODO(geoah) It makes the assumption that the payload is JSON.
// TODO(geoah) This doesn't seem to be used any more.
type EventUnknown struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

// Repository is responsible of getting events from a journal and
// aggregating them in `Aggregate`s based on the events' `guid`s.
// Currently it holds them in key/value pairs in-memory, but should
// accept a store object that should handle aggregate get/set operations.
type Repository struct {
	aggregate Aggregate
	event     Event
	pairs     map[string]Aggregate
	store     store.Store
}

// NewRepository creates a new `Repository` given a store, `Aggregate` template,
// and `Event` template. The two templates should be just empty structs that will
// be used as "guides" to create the new aggregates.
func NewRepository(store store.Store, aggregate Aggregate, event Event) *Repository {
	return &Repository{
		aggregate: aggregate,
		event:     event,
		store:     store,
		pairs:     map[string]Aggregate{},
	}
}

// GetByGUID returns an `Aggregate` from it's `guid`, if it exists.
// Else if will return `Err`
func (r *Repository) GetByGUID(key []byte) (Aggregate, error) {
	if a, ok := r.pairs[string(key)]; ok {
		return a, nil
	}
	return nil, ErrAggregateNotFound
}

// AppendedEntry satisfies the `journal.Notifiee` interface and is called when
// a new entry has been added to the `Journal`.
func (r *Repository) AppendedEntry(entry journal.Entry) {
	fmt.Println("> Processing", string(entry.GetIndex()), string(entry.GetPayload()))
	// TODO(geoah) Check that this event hasn't already been processed.
	// TODO(geoah) Check that the previous event has already been processed.

	// create a new instance of our aggregate from our template
	aggregatePtr := reflect.New(reflect.TypeOf(r.aggregate).Elem())
	aggregate := aggregatePtr.Interface().(Aggregate)

	// create a new instance of our event from our template
	eventPtr := reflect.New(reflect.TypeOf(r.aggregate).Elem())
	event := eventPtr.Interface().(Event)

	// decode the event
	event.Unmarshal(entry.GetPayload()) // TODO(geoah) handle error

	// check if there we already have an aggregate with the same guid
	exists := false
	if a, ok := r.pairs[string(event.GetGuid())]; ok {
		exists = true
		aggregate = a
	}

	// apply the event on our aggregate
	aggregate.ApplyEvent(event)

	// if the aggregate is new let's store
	if exists == false {
		r.pairs[string(event.GetGuid())] = aggregate
	}

	pretty.Println("r.pairs", r.pairs)
}
