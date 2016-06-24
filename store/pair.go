package store

// Key for use in the ClusteringKey.
type Key interface{}

// ClusteringKey is a composition of one of more Keys.
// A ClusteringKey can be either complete or incomplete.
// Complete ClusteringKeys can be persisted, while non-complete ones can
// only be used for querying.
// eg. `func (ck *ClusteringKey) IsComplete() { return len(ck.GetKeys()) == 3 }`
type ClusteringKey interface {
	// GetKeys returns a list of the Keys that comprise the ClusteringKey.
	GetKeys() []Key
	// IsComplete checks the ClusteringKey for completenes.
	IsComplete() bool
}

// Value of the key-value pairs
type Value interface{}
