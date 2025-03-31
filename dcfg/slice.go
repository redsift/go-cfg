package dcfg

import (
	"context"
	"errors"
	"fmt"
	"slices"
)

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
func NewTypedSlice[T any](b Backend, key Key) (*TypedSlice[T], error) {
	type_ := TypeOf[[]T]()
	if key.Type != type_ {
		return nil, fmt.Errorf("invalid type %q in Slice key, expected %q", key.Type, type_)
	}
	return &TypedSlice[T]{
		backend: b,
		key:     key,
		slice:   b.Slice(key),
	}, nil
}

// TypedSlice wraps a Slice implementation from a Backend and provides type-safety.
type TypedSlice[T any] struct {
	backend Backend
	key     Key
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
	err := t.slice.Load(ctx, &out)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
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

// Subscribe subscribes to updates on the store
func (t *TypedSlice[T]) Subscribe(ctx context.Context, update func(updated []T, err error) bool) error {
	return Subscribe[[]T](ctx, t.backend, t.key, func(updated []T, m Meta, err error) bool {
		return update(updated, err)
	})
}

// SubscribeDiff subscribes to updates on the store and calculates the diff
func (t *TypedSlice[T]) SubscribeDiff(ctx context.Context, compare func(T, T) int, update func(add, remove []T, err error) bool) error {
	cur, err := t.Load(ctx)
	if err != nil {
		return nil
	}

	slices.SortFunc(cur, compare)
	cur = slices.CompactFunc(cur, func(a, b T) bool {
		return compare(a, b) == 0
	})

	return Subscribe[[]T](ctx, t.backend, t.key, func(updated []T, m Meta, err error) bool {
		if err != nil {
			return update(nil, nil, err)
		}

		slices.SortFunc(updated, compare)
		updated = slices.CompactFunc(updated, func(a, b T) bool {
			return compare(a, b) == 0
		})

		var (
			added, removed []T
			c, u           int
		)

		for c < len(cur) && u < len(updated) {
			diff := compare(cur[c], updated[u])
			if diff == 0 {
				c++
				u++
			} else if diff > 0 {
				added = append(added, updated[u])
				u++
			} else { // diff < 0
				removed = append(removed, cur[c])
				c++
			}
		}

		if c < len(cur) {
			removed = append(removed, cur[c:]...)
		}

		if u < len(updated) {
			added = append(added, updated[u:]...)
		}

		cur = updated

		return update(added, removed, nil)
	})
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
