package store

import "errors"

// Key of the key-valye pairs.
type Key []byte

// Value of the key-value pairs.
type Value []byte

// ErrKeyNotFound when the value for the key provided could not be retrieved.
var ErrKeyNotFound = errors.New("Key not found.")

// ErrInternalError when something doesn't go as expected.
var ErrInternalError = errors.New("Internal error.")

// Store is a very generic interface for storing key-value pairs.
type Store interface {
	// Put sets the key's value, overwriting the previous if it exists,
	// errors with `ErrInternalError`.
	Put(key Key, value Value) (err error)
	// Get gets the value for a given key,else errors with
	// `ErrKeyNotFound`, or `ErrInternalError`.
	Get(key Key) (value Value, err error)
	// Delete removes a value if it exists, else errors with
	// `ErrKeyNotFound`, or `ErrInternalError`.
	Delete(key Key) (err error)
}
