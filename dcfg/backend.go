package dcfg

import (
	"context"
)

// Meta currently only contains the value generation used for optimistic locking
type Meta struct {
	Generation uint
}

// Backend must be implemented by a configuration store
type Backend interface {
	// Delete removes the given `key` from the store.
	Delete(ctx context.Context, key Key) error

	// Load loads the given `key` from the store into `target`.
	Load(ctx context.Context, key Key, target any) (Meta, error)

	// Store writes the given value with the given `key` to the store.
	// if meta is nil, the value will always be overwritten.
	// if meta is not nil, the value will only be written if the generation in the store matches the given value.
	Store(ctx context.Context, key Key, meta *Meta, value any) error

	// Subscribe creates a Stream that notifies on changes to the given `key`.
	Subscribe(ctx context.Context, key Key) (Stream, error)

	// Map creates and returns a Map implementation for the given `key`.
	Map(key Key) Map

	// Slice creates and returns a Slice implementation for the given `key`.
	Slice(key Key) Slice
}
