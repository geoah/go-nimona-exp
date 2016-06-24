package store

import "errors"

// ErrClusteringKeyNotFound when a complete clustering key was provided
// but could not be retrieved.
var ErrClusteringKeyNotFound = errors.New("ClusteringKey not found.")

// ErrClusteringKeyNotComplete when a clustering key is not complete, so
// it cannot be persisted.
var ErrClusteringKeyNotComplete = errors.New("ClusteringKey is not complete.")

// ErrInternalError when something doesn't go as expected.
var ErrInternalError = errors.New("Internal error.")

// Store is a very generic interface for storing key-value pairs.
type Store interface {
	// Put sets the key's value, overwriting the previous if it exists.
	Put(completeKey ClusteringKey, value Value) (err error)
	// GetOne gets the value for a clustering key and updates theresult, else
	// errors with `ErrClusteringKeyNotFound`, or `ErrClusteringKeyNotComplete`.
	GetOne(completeKey ClusteringKey, result Value) (err error)
	// GetAll updates the results list with the values of the given incomplete
	// ClusteringKey, or `ErrClusteringKeyComplete`
	GetAll(key ClusteringKey, results []*Value) (err error)
	// Delete removed the key's value if it exists, else errors with
	// `ErrClusteringKeyNotFound`, or `ErrClusteringKeyNotComplete`.
	Delete(completeKey ClusteringKey) (err error)
}
