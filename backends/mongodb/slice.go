package mongodb

import (
	"context"

	"github.com/redsift/go-cfg/dcfg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Slice implements dcfg.Backend.
func (b *Backend) Slice(key dcfg.Key) dcfg.Slice {
	return &Slice{backend: b, key: key}
}

// Slice implements dcfg.Slice
type Slice struct {
	backend *Backend
	key     dcfg.Key
}

// Append implements dcfg.Slice.
func (s *Slice) Append(ctx context.Context, items ...any) error {
	if len(items) == 0 {
		return nil
	}

	ops := make(bson.D, len(items))
	for i, item := range items {
		ops[i] = bson.E{Key: "value", Value: item}
	}

	return s.backend.withColl(ctx, s.key, func(coll *mongo.Collection) error {
		_, err := coll.UpdateOne(
			ctx,
			s.backend.filter(s.key),
			bson.E{Key: "$push", Value: ops},
		)
		return err
	})
}

// RemoveIndexes implements dcfg.Slice.
func (s *Slice) Load(ctx context.Context, target any) error {
	return s.backend.Load(ctx, s.key, target)
}

// RemoveIndexes implements dcfg.Slice.
func (s *Slice) RemoveItems(ctx context.Context, items ...any) error {
	if len(items) == 0 {
		return nil
	}

	ops := make(bson.D, len(items))
	for i, item := range items {
		ops[i] = bson.E{Key: "value", Value: item}
	}

	return s.backend.withColl(ctx, s.key, func(coll *mongo.Collection) error {
		_, err := coll.UpdateOne(
			ctx,
			s.backend.filter(s.key),
			bson.E{Key: "$pull", Value: ops},
		)
		return err
	})
}

// Store implements dcfg.Slice.
func (s *Slice) Store(ctx context.Context, items ...any) error {
	return s.backend.Store(ctx, s.key, items)
}
