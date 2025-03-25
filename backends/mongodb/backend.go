package mongodb

import (
	"context"
	"strings"
	"sync"

	"github.com/redsift/go-cfg/dcfg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func New(client *mongo.Client) dcfg.Backend {
	return &Backend{
		client:      client,
		collections: map[string]*mongo.Collection{},
	}
}

type Backend struct {
	client      *mongo.Client
	lock        sync.Mutex
	collections map[string]*mongo.Collection
}

// Load implements dcfg.Backend.
func (b *Backend) Load(ctx context.Context, key dcfg.Key, target any) error {
	var tmp envelope[bson.Raw]
	if err := b.withColl(ctx, key, func(coll *mongo.Collection) error {
		return coll.FindOne(ctx, b.filter(key)).Decode(&tmp)
	}); err != nil {
		return err
	}
	return bson.Unmarshal(tmp.Value, target)
}

// Store implements dcfg.Backend.
func (b *Backend) Store(ctx context.Context, key dcfg.Key, value any) error {
	tmp := envelope[any]{
		Key:     strings.Join(key.Elements, "/"),
		Version: uint(key.Version),
		Value:   value,
	}
	return b.withColl(ctx, key, func(coll *mongo.Collection) error {
		_, err := coll.ReplaceOne(ctx, tmp, options.Replace().SetUpsert(true))
		return err
	})
}
