package nilbe

import (
	"context"

	"github.com/redsift/go-cfg/dcfg"
)

var Nil dcfg.Backend = Backend{}

type Backend struct{}

// Delete implements dcfg.Backend.
func (b Backend) Delete(ctx context.Context, key dcfg.Key) error {
	return nil
}

// Load implements dcfg.Backend.
func (b Backend) Load(ctx context.Context, key dcfg.Key, target any) (dcfg.Meta, error) {
	return dcfg.Meta{}, dcfg.ErrNotFound
}

// Map implements dcfg.Backend.
func (b Backend) Map(key dcfg.Key) dcfg.Map {
	return Map{}
}

// Slice implements dcfg.Backend.
func (b Backend) Slice(key dcfg.Key) dcfg.Slice {
	return Slice{}
}

// Store implements dcfg.Backend.
func (b Backend) Store(ctx context.Context, key dcfg.Key, meta *dcfg.Meta, value any) error {
	return nil
}

// Subscribe implements dcfg.Backend.
func (b Backend) Subscribe(ctx context.Context, key dcfg.Key) (dcfg.Stream, error) {
	return Stream{}, nil
}

var _ dcfg.Map = Map{}

type Map struct{}

// Clear implements dcfg.Map.
func (m Map) Clear(ctx context.Context) error {
	return nil
}

// DelKey implements dcfg.Map.
func (m Map) DelKey(ctx context.Context, key string) error {
	return nil
}

// GetKey implements dcfg.Map.
func (m Map) GetKey(ctx context.Context, key string, target any) error {
	return dcfg.ErrNotFound
}

// Load implements dcfg.Map.
func (m Map) Load(ctx context.Context, target any) error {
	// no value -> empty map
	return nil
}

// SetKey implements dcfg.Map.
func (m Map) SetKey(ctx context.Context, key string, value any) error {
	return nil
}

// Update implements dcfg.Map.
func (m Map) Update(ctx context.Context, values map[string]any) error {
	return nil
}

type Slice struct{}

// Append implements dcfg.Slice.
func (s Slice) Append(ctx context.Context, items ...any) error {
	return nil
}

// Load implements dcfg.Slice.
func (s Slice) Load(ctx context.Context, items any) error {
	return dcfg.ErrNotFound
}

// RemoveItems implements dcfg.Slice.
func (s Slice) RemoveItems(ctx context.Context, items ...any) error {
	return nil
}

// Store implements dcfg.Slice.
func (s Slice) Store(ctx context.Context, items ...any) error {
	return nil
}

type Stream struct{}

// Close implements dcfg.Stream.
func (s Stream) Close() error {
	return nil
}

// Decode implements dcfg.Stream.
func (s Stream) Decode(any) (dcfg.Meta, error) {
	return dcfg.Meta{}, nil
}

// Next implements dcfg.Stream.
func (s Stream) Next(context.Context) bool {
	return false
}
