package dcfg

import (
	"context"
)

// Backend must be implemented by a configuration store
type Backend interface {
	// Load loads the given `key` from the store into `target`.
	Load(ctx context.Context, key Key, target any) error

	// Store writes the given value with the given `key` to the store.
	Store(ctx context.Context, key Key, value any) error

	// Subscribe creates a Stream that notifies on changes to the given `key`.
	Subscribe(ctx context.Context, key Key) (Stream, error)

	// Slice creates and returns a Slice implementation for the given `key`.
	Slice(key Key) Slice
}
