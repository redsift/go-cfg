package dcfg

import (
	"context"
	"fmt"
)

// Map defines an interface to interact with a string value map.
type Map interface {
	// DelKey removes a key from the map.
	DelKey(ctx context.Context, key string) error
	// GetKey retrieves a single value from the map.
	GetKey(ctx context.Context, key string, target any) error
	// SetKey inserts (or overwrites) the given value at the given key.
	SetKey(ctx context.Context, key string, value any) error

	// Clear removes all map entries
	Clear(ctx context.Context) error
	// Load loads all key/value pairs into the given target. The target must be a map with a ~string key type.
	Load(ctx context.Context, target any) error
	// Update overwrites the whole map with the given values.
	Update(ctx context.Context, values map[string]any) error
}

// NewTypedMap creates a new TypedMap.
func NewTypedMap[K ~string, V any](b Backend, key Key) (*TypedMap[K, V], error) {
	type_ := TypeOf[map[K]V]()
	if key.Type != type_ {
		return nil, fmt.Errorf("invalid type %q in Map key, expected %q", key.Type, type_)
	}
	return &TypedMap[K, V]{
		backend: b,
		key:     key,
		m:       b.Map(key),
	}, nil
}

// TypedMap wraps a Map to provide strong type safety.
type TypedMap[K ~string, V any] struct {
	backend Backend
	key     Key
	m       Map
}

func (t *TypedMap[K, V]) DelKey(ctx context.Context, key K) error {
	return t.m.DelKey(ctx, string(key))
}

func (t *TypedMap[K, V]) GetKey(ctx context.Context, key K) (value V, err error) {
	err = t.m.GetKey(ctx, string(key), &value)
	return
}

func (t *TypedMap[K, V]) SetKey(ctx context.Context, key K, value V) error {
	return t.m.SetKey(ctx, string(key), value)
}

func (t *TypedMap[K, V]) Clear(ctx context.Context) error {
	return t.m.Clear(ctx)
}

func (t *TypedMap[K, V]) Load(ctx context.Context) (out map[K]V, err error) {
	err = t.m.Load(ctx, &out)
	return
}

func (t *TypedMap[K, V]) Update(ctx context.Context, values map[K]V) error {
	tmp := make(map[string]any, len(values))
	for k, v := range values {
		tmp[string(k)] = v
	}
	return t.m.Update(ctx, tmp)
}
