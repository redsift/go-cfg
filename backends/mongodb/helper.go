package mongodb

import (
	"context"
	"regexp"
	"strings"

	"github.com/redsift/go-cfg/dcfg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var mangleRE = regexp.MustCompile("[^a-zA-Z0-9_]")

type envelope[Value any] struct {
	Key     string `json:"key"`
	Version uint   `json:"version"`
	Value   Value  `json:"value"`
}

func (b *Backend) coll(ctx context.Context, key dcfg.Key) (*mongo.Collection, error) {
	collName := key.Elements[0] + "_" + mangleRE.ReplaceAllString(string(key.Type), "_")

	b.lock.Lock()
	defer b.lock.Unlock()

	coll, ok := b.collections[collName]
	if ok {
		return coll, nil
	}

	coll = b.client.Database("dcfg").Collection(key.Elements[0] + "_")
	indexes := coll.Indexes()
	_, err := indexes.CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "key", Value: 1},
			{Key: "version", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	return coll, err
}

func (b *Backend) filter(key dcfg.Key) bson.D {
	return bson.D{
		{Key: "key", Value: strings.Join(key.Elements, "/")},
		{Key: "value", Value: uint(key.Version)},
	}
}

func (b *Backend) withColl(ctx context.Context, key dcfg.Key, fn func(*mongo.Collection) error) error {
	coll, err := b.coll(ctx, key)
	if err != nil {
		return err
	}
	return fn(coll)
}
