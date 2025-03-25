package dcfg

import (
	"context"
)

// Backend must be implemented by a configuration store
type Backend interface {
	// Load loads the given `type_` with `key` from the store into `target`.
	Load(ctx context.Context, type_ Type, key Key, target any) error
	// Store writes the given value with `type_` and `key` to the store.
	Store(ctx context.Context, type_ Type, key Key, value any) error
	// Subscribe creates a Stream that notifies on changes to the given `type_` and `key`.
	Subscribe(ctx context.Context, type_ Type, key Key) (Stream, error)

	Slice(type_ Type, key Key) Slice
}
