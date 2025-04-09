package mongodb

import (
	"context"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/redsift/go-cfg/dcfg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// collectionMangleRE matches all characters that are not safe to use in a collection name
var collectionMangleRE = regexp.MustCompile("[^a-zA-Z0-9_]")

// envelope is used to wrap the value with the key fields.
type envelope[Value any] struct {
	Key        string `json:"key"`
	Version    uint   `json:"version"`
	Value      Value  `json:"value"`
	Generation uint   `json:"generation"`
}

// coll derives a collection name from the key and ensures the collection is set up with the unique
// index.
func (b *Backend) coll(ctx context.Context, key dcfg.Key) (*mongo.Collection, error) {
	collName := key.Elements[0] + "_" + collectionMangleRE.ReplaceAllString(string(key.Type), "_")

	b.lock.Lock()
	defer b.lock.Unlock()

	coll, ok := b.collections[collName]
	if ok {
		return coll, nil
	}

	coll = b.client.Database(b.dbName).Collection(key.Elements[0] + "_")
	indexes := coll.Indexes()
	_, err := indexes.CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "key", Value: 1},
			{Key: "version", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	return coll, mapError(err)
}

// filter returns the filter for the key
func (b *Backend) filter(key dcfg.Key) bson.D {
	return bson.D{
		{Key: "key", Value: strings.Join(key.Elements, "/")},
		{Key: "version", Value: uint(key.Version)},
	}
}

// withColl calls the given func with the collection derived from the given key.
func (b *Backend) withColl(ctx context.Context, key dcfg.Key, fn func(*mongo.Collection) error) error {
	coll, err := b.coll(ctx, key)
	if err != nil {
		return err
	}
	return mapError(fn(coll))
}

// mapError maps mongo errors to dcfg errors
func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return dcfg.ErrNotFound
	}
	if errors.Is(err, io.EOF) {
		return dcfg.ErrNotFound
	}
	return err
}
