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
	// Put sets the key's value, overwriting the previous if it exists.
	Put(key Key, value Value) (err error)
	// GetOne gets the value for a key and updates theresult, else
	// errors with `ErrKeyNotFound`.
	GetOne(key Key) (value Value, err error)
	// GetAll finds all pairs that partially match the given key
	// (left to right).
	GetAll(key Key) (results []*Value, err error)
	// Delete removes the key's value if it exists, else errors with
	// `ErrKeyNotFound`.
	Delete(key Key) (err error)
}
