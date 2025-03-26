package mongodb

import (
	"context"

	"github.com/redsift/go-cfg/dcfg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	op := bson.D{{
		Key: "$push", Value: bson.D{{
			Key: "value", Value: bson.D{{
				Key: "$each", Value: items,
			}},
		}},
	}}

	return s.backend.withColl(ctx, s.key, func(coll *mongo.Collection) error {
		_, err := coll.UpdateOne(ctx, s.backend.filter(s.key), op, options.UpdateOne().SetUpsert(true))
		return err
	})
}

// RemoveIndexes implements dcfg.Slice.
func (s *Slice) Load(ctx context.Context, target any) error {
	//var tmp bson.Raw
	var tmp envelope[bson.Raw]
	if err := s.backend.load(ctx, s.key, &tmp); err != nil {
		return err
	}
	return bson.UnmarshalValue(bson.TypeArray, tmp.Value, target)
}

// RemoveIndexes implements dcfg.Slice.
func (s *Slice) RemoveItems(ctx context.Context, items ...any) error {
	if len(items) == 0 {
		return nil
	}

	op := bson.D{{
		Key: "$pull", Value: bson.D{{
			Key: "value", Value: bson.D{{
				Key: "$in", Value: items,
			}},
		}},
	}}

	return s.backend.withColl(ctx, s.key, func(coll *mongo.Collection) error {
		_, err := coll.UpdateOne(ctx, s.backend.filter(s.key), op)
		return err
	})
}

// Store implements dcfg.Slice.
func (s *Slice) Store(ctx context.Context, items ...any) error {
	return s.backend.Store(ctx, s.key, items)
}
