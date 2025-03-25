package dcfg

import "context"

// Slice abstracts a list of elements
type Slice interface {
	// Append adds an item to the slice
	Append(ctx context.Context, items ...any) error
	// Load loads the whole slice
	Load(ctx context.Context, items any) error
	// Store overwrites all values in a slice
	Store(ctx context.Context, items ...any) error
	// RemoveItems removes the given items from the slice
	RemoveItems(ctx context.Context, items ...any) error
}

// NewTypedSlice creates a new TypedSlice
func NewTypedSlice[T any](b Backend, key Key) *TypedSlice[T] {
	type_ := TypeOf([]T{})
	return &TypedSlice[T]{
		backend: b,
		key:     key,
		type_:   type_,
		slice:   b.Slice(type_, key),
	}
}

// TypedSlice wraps a Slice implementation from a Backend and provides type-safety.
type TypedSlice[T any] struct {
	backend Backend
	key     Key
	type_   Type
	slice   Slice
}

// Append adds an item to the slice
func (t *TypedSlice[T]) Append(ctx context.Context, items ...T) error {
	if len(items) == 0 {
		return nil
	}

	a := make([]any, len(items))
	for i, e := range items {
		a[i] = e
	}

	return t.slice.Append(ctx, a...)
}

// Load loads the whole slice
func (t *TypedSlice[T]) Load(ctx context.Context) (out []T, _ error) {
	err := t.backend.Load(ctx, t.type_, t.key, &out)
	return out, err
}

// Store overwrites all values in a slice
func (t *TypedSlice[T]) Store(ctx context.Context, items ...T) error {
	a := make([]any, len(items))
	for i, e := range items {
		a[i] = e
	}

	return t.slice.Store(ctx, a...)
}

// RemoveItems removes the given items from the slice
func (t *TypedSlice[T]) RemoveItems(ctx context.Context, items ...T) error {
	if len(items) == 0 {
		return nil
	}

	a := make([]any, len(items))
	for i, e := range items {
		a[i] = e
	}

	return t.slice.RemoveItems(ctx, a...)
}
