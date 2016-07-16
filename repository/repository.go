package repository

import "github.com/nimona/go-nimona/journal"

// Resource is the value of any resource.
type Resource interface{}

// Repository handles the events from a single stream.
// It is responsible of managing the projections (resources) for the various event types.
type Repository interface {
	GetResourceByID(string) (Resource, error)
	AppendEntry(journal.Entry) (Resource, error)
	AppendedEntry(journal.Entry)
}
